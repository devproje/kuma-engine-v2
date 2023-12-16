package net.projecttl.kuma.engine.command

import kotlinx.coroutines.coroutineScope
import kotlinx.coroutines.launch
import net.dv8tion.jda.api.JDA
import net.dv8tion.jda.api.events.interaction.command.SlashCommandInteractionEvent
import net.dv8tion.jda.api.hooks.ListenerAdapter
import org.slf4j.Logger

open class CommandHandler(
    val guildId: String? = null
) : ListenerAdapter() {
    private val commands = mutableListOf<CommandExecutor>()

    fun addCommands(vararg commands: CommandExecutor) {
        this.commands.addAll(commands)
    }

    fun dropCommands(command: CommandExecutor) {
        if (!commands.contains(command)) {
            return
        }

        commands.remove(command)
    }

    override fun onSlashCommandInteraction(event: SlashCommandInteractionEvent) {
        val command = commands.find { it.data.name == event.name }
        if (command == null) {
            return
        }

        try {
            command.execute(event)
        } catch (ex: Exception) {
            event.reply("Error occurred while executing command.").queue()
            ex.printStackTrace()
        }
    }

    suspend fun register(jda: JDA, logger: Logger) {
        coroutineScope {
            launch {
                for (cmd in commands) {
                    if (guildId == null) {
                        jda.upsertCommand(cmd.data.build()).queue()
                        logger.info("Registered global command: /${cmd.data.name}")

                        continue
                    }

                    val guild = jda.getGuildById(guildId) ?: throw Exception("current guild id is not exist")

                    guild.upsertCommand(cmd.data.build()).queue()
                    logger.info("Registered ${guild.id} command: /${cmd.data.name}")
                }
            }
        }
    }
}