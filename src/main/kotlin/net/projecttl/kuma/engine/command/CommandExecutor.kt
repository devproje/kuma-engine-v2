package net.projecttl.kuma.engine.command

import net.dv8tion.jda.api.events.interaction.command.SlashCommandInteractionEvent
import net.dv8tion.jda.api.interactions.commands.build.CommandData

interface CommandExecutor {
    val data: CommandData
    fun executor(event: SlashCommandInteractionEvent)
}