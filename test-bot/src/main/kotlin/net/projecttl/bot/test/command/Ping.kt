package net.projecttl.bot.test.command

import net.dv8tion.jda.api.events.interaction.command.SlashCommandInteractionEvent
import net.projecttl.kuma.engine.model.CommandExecutor
import net.projecttl.kuma.engine.build.CommandDataBuilder
import net.projecttl.kuma.engine.build.EmojiBuilder

object Ping : CommandExecutor {

    override val data = CommandDataBuilder().apply {
        name = "ping"
        description = "Getting discord API latency speed"
    }

    override fun execute(event: SlashCommandInteractionEvent) {
        val emoji = EmojiBuilder("ping_pong").build()
        event.reply("**$emoji Pong!** ${event.jda.gatewayPing}**ms**").queue()
    }
}