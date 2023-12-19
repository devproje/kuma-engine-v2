package net.projecttl.bot.test

import net.dv8tion.jda.api.events.session.ReadyEvent
import net.dv8tion.jda.api.hooks.ListenerAdapter
import net.projecttl.bot.test.command.Ping
import net.projecttl.bot.test.command.TestOwner
import net.projecttl.kuma.engine.KumaEngine
import org.slf4j.LoggerFactory

private val logger = LoggerFactory.getLogger("MainKt")

suspend fun main() {
    val bot = KumaEngine(Config.token)
    bot.command.addCommands(Ping, TestOwner)
    bot.addHandler(Ready)

    bot.build()
}

object Ready : ListenerAdapter() {
    override fun onReady(event: ReadyEvent) {
        logger.info("Logged in as ${event.jda.selfUser.name}")
    }
}
