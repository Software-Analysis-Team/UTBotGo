#include <llvm/IR/Module.h>
#include <llvm/IR/LLVMContext.h>
#include <llvm/IR/Function.h>
#include <llvm/IR/BasicBlock.h>
#include <llvm/IR/Instruction.h>
#include <llvm/IR/Instructions.h>
#include <llvm/IRReader/IRReader.h>
#include <llvm/Support/SourceMgr.h>
#include <llvm/Support/raw_ostream.h>
#include <llvm/Bitcode/BitcodeWriter.h>

#include <iostream>
#include <memory>
#include <utility>

using namespace std;
using namespace llvm;

int main(int argc, char *argv[])
{
    LLVMContext context;
    SMDiagnostic error;
    unique_ptr<Module> module = parseIRFile(argv[1], error, context);
    auto voidT = Type::getVoidTy(context);
    auto i8T = Type::getInt8Ty(context);
    auto i8PT = Type::getInt8PtrTy(context);
    auto i64T = Type::getInt64Ty(context);
    for (auto& glob : module->globals()) {
        if (glob.isExternallyInitialized()) {
            glob.setExternallyInitialized(false);
            glob.setInitializer(ConstantAggregateZero::get(glob.getType()->getPointerElementType()));
        }
    }
    auto kleeMakeSymbolic = module->getFunction("klee_make_symbolic");
    if (! kleeMakeSymbolic) {
        Twine name = "klee_make_symbolic";
        kleeMakeSymbolic = Function::Create(FunctionType::get(voidT, {i8PT, i64T, i8PT}, false), GlobalValue::ExternalLinkage, name, *module);
    }
    if (auto funcKleeMain = module->getFunction("klee.main")) {
        funcKleeMain->eraseFromParent();
    }
    Twine name = "klee.main";
    Twine emptyTwine = Twine();
    auto kleeMain = Function::Create(FunctionType::get(voidT, false), GlobalValue::InternalLinkage, name, *module);
    if (auto func = module->getFunction(StringRef(argv[3]))) {
        Twine entryBBname = "entry";
        auto entryBB = BasicBlock::Create(context, entryBBname, kleeMain);
        int argNum = 0;
        auto kleeZero = new GlobalVariable(*module, i64T, true, GlobalValue::InternalLinkage, ConstantInt::get(i64T, 0), "klee.zero");
        vector<Value*> args(func->arg_size());
        for (auto arg = func->arg_begin(); arg != func->arg_end(); arg++, argNum++) {
            if (argNum != 0) {
                auto argName = arg->getName().str();
                auto paramName = Twine("klee." + argName + ".addr");
                auto paramType = arg->getType();
                auto addr = new AllocaInst(paramType, 0, paramName, entryBB);
                int paramNameArraySize = argName.size() + 1;
                vector<Constant*> paramNameArray(paramNameArraySize);
                for (int i = 0; i < paramNameArraySize - 1; i++) {
                    paramNameArray[i] = ConstantInt::get(i8T, argName[i]);
                }
                paramNameArray[paramNameArraySize-1] = ConstantInt::get(i8T, 0);
                ArrayRef<Constant*> paramNameArrayRef(paramNameArray);
                auto paramNameType = ArrayType::get(i8T, paramNameArraySize);
                auto paramNameValue = ConstantArray::get(paramNameType, paramNameArrayRef);
                auto paramNameName = Twine("klee." + argName + ".str");
                auto gv = new GlobalVariable(*module, paramNameType, true, GlobalValue::InternalLinkage, paramNameValue, paramNameName);
                auto paramBCname = Twine("klee." + argName + ".bc");
                auto bc = new BitCastInst(addr, i8PT, paramBCname, entryBB);
                auto paramTypeSize = paramType->getPrimitiveSizeInBits().getKnownMinSize() / 8;
                auto paramNameGEPIname = Twine("klee." + argName + ".name");
                auto paramCIname = Twine("klee." + argName + ".kms");
                auto gepiVar = new GlobalVariable(*module, i8T, true, GlobalVariable::InternalLinkage, ConstantInt::get(i8T, paramNameArraySize), "klee_modify.." + gv->getName().str());
                CallInst::Create(kleeMakeSymbolic->getFunctionType(), kleeMakeSymbolic, {bc, ConstantInt::get(i64T, paramTypeSize), gepiVar}, "", entryBB);
                auto ld = new LoadInst(addr->getType()->getPointerElementType(), addr, "klee." + argName + ".ld", entryBB);
                args[argNum] = ld;
            } else {
                args[argNum] = UndefValue::get(i8PT);
            }
        }
        CallInst::Create(func->getFunctionType(), func, args, "klee.call", entryBB);
        ReturnInst::Create(context, entryBB);
    }
    error_code ec;
    raw_fd_ostream out(argv[2], ec);
    module->print(out, nullptr);
    out.flush();
}
