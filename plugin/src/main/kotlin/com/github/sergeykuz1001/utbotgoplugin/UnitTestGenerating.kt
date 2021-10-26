package com.github.sergeykuz1001.utbotgoplugin

import com.intellij.openapi.actionSystem.AnAction
import com.intellij.openapi.actionSystem.AnActionEvent
import com.intellij.openapi.actionSystem.LangDataKeys
import com.intellij.openapi.ui.Messages
import java.io.BufferedReader
import java.io.InputStreamReader
import java.util.*

class UnitTestGenerating : AnAction() {
    private fun showErrorMessage(message: String) {
        Messages.showMessageDialog(message, "Error", Messages.getErrorIcon())
    }

    override fun actionPerformed(event: AnActionEvent) {
        val file = event.getData(LangDataKeys.PSI_FILE)
        if (file == null) {
            showErrorMessage("Can't find out selected file")
            return
        }
        val path = file.originalFile.virtualFile.path
        val proc = Runtime.getRuntime().exec("utbotgo --gofuzz $path")
        val scan = Scanner(BufferedReader(InputStreamReader(proc.inputStream)))
        val genFilesAmt = scan.nextInt()
        scan.close()
        if (genFilesAmt == 0) showErrorMessage("Nothing files were generated")
        val exitVal = proc.waitFor()
        if (exitVal != 0) showErrorMessage("Exit code of main program is $exitVal")
    }
}
