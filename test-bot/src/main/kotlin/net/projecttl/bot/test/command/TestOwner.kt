package net.projecttl.bot.test.command

import net.dv8tion.jda.api.events.interaction.command.SlashCommandInteractionEvent
import net.projecttl.bot.test.Config
import net.projecttl.kuma.engine.model.CommandExecutor
import net.projecttl.kuma.engine.build.CommandDataBuilder

object TestOwner : CommandExecutor {
    override val data = CommandDataBuilder().apply {
        name = "owner"
        description = "test"
    }

    override fun execute(event: SlashCommandInteractionEvent) {
        if (event.user.id != Config.owner) {
            event.reply("You're not bot owner").queue()
            return
        }

        event.reply("Hello! ${event.user.asMention}").queue()
    }
}