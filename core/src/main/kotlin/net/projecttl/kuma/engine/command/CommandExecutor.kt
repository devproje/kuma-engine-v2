package net.projecttl.kuma.engine.command

import net.dv8tion.jda.api.events.interaction.command.SlashCommandInteractionEvent
import net.projecttl.kuma.engine.`object`.CommandDataBuilder

interface CommandExecutor {
    val data: CommandDataBuilder
    fun execute(event: SlashCommandInteractionEvent)
}