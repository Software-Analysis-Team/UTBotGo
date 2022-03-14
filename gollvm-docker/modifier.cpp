#include <llvm/IR/Module.h>
#include <llvm/IR/LLVMContext.h>
#include "llvm/Support/SourceMgr.h"
#include <llvm/IRReader/IRReader.h>

using namespace std;
using namespace llvm;

int main(int argc, char *argv[])
{
    LLVMContext context;
    SMDiagnostic error;
    unique_ptr<Module> module = parseIRFile(argv[1], error, context);
    for (auto func = module->begin(); func != module->end(); func++) {
        if (func->getName() == "main.main") {
            func->setName("main");
        }
    }
    module->dump();
}
