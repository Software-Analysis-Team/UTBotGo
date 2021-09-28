package com.github.sergeykuz1001.utbotgoplugin

import com.intellij.openapi.actionSystem.AnAction
import com.intellij.openapi.actionSystem.AnActionEvent
import com.intellij.openapi.actionSystem.LangDataKeys
import com.intellij.openapi.ui.Messages

class UnitTestGenerating : AnAction() {
    override fun actionPerformed(event: AnActionEvent) {
        val file = event.getData(LangDataKeys.PSI_FILE)
        if (file == null) {
            Messages.showMessageDialog("Can't find out selected file", "Problem", Messages.getErrorIcon())
            return
        }
        val path = file.originalFile.virtualFile.path
        val answer = Messages.showDialog(
            "Do you want to generate unit tests for \"$path\"?",
            "Generating Unit Tests",
            arrayOf("Yes", "No"),
            0,
            Messages.getQuestionIcon(),
        )
        if (answer != 0) return
        Runtime.getRuntime().exec("UTBotGo $path")
    }
}