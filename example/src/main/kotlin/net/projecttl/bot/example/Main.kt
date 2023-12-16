package org.example.net.projecttl.bot.example

import net.dv8tion.jda.api.events.session.ReadyEvent
import net.dv8tion.jda.api.hooks.ListenerAdapter
import net.projecttl.kuma.engine.KumaEngine

suspend fun main() {
    val bot = KumaEngine(Config.token)
//    val cmd = CommandHandler()
//    cmd.addCommands()

//    bot.addCommandHandler(cmd)
    bot.addHandler(Ready)

    bot.build()
}

object Ready : ListenerAdapter() {
    override fun onReady(event: ReadyEvent) {
        println("Logged in as ${event.jda.selfUser.name}")
    }
}
