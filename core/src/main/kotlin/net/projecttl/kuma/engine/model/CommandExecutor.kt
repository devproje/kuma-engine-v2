package net.projecttl.kuma.engine.model

import net.dv8tion.jda.api.events.interaction.command.SlashCommandInteractionEvent
import net.projecttl.kuma.engine.build.CommandDataBuilder

interface CommandExecutor {
    val data: CommandDataBuilder
    fun execute(event: SlashCommandInteractionEvent)
}