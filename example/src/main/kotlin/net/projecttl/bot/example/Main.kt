package net.projecttl.bot.example

import net.dv8tion.jda.api.events.session.ReadyEvent
import net.dv8tion.jda.api.hooks.ListenerAdapter
import net.projecttl.kuma.engine.KumaEngine
import org.slf4j.LoggerFactory

private val logger = LoggerFactory.getLogger("net.projecttl.bot.example.MainKt")

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
        logger.info("Logged in as ${event.jda.selfUser.name}")
    }
}
