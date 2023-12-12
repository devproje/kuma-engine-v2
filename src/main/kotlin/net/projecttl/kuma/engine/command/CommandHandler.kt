package net.projecttl.kuma.engine.command

import kotlinx.coroutines.coroutineScope
import kotlinx.coroutines.launch
import net.dv8tion.jda.api.JDA
import net.dv8tion.jda.api.events.interaction.command.SlashCommandInteractionEvent
import net.dv8tion.jda.api.hooks.ListenerAdapter
import net.projecttl.kuma.engine.enum.Lang
import java.util.logging.Logger

class CommandHandler(
    val name: String = "default",
    val guildId: String? = null,
    val lang: Lang = Lang.EN_US
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
            command.executor(event)
        } catch (ex: Exception) {
            event.reply("Error occurred while executing command.").queue()
            ex.printStackTrace()
        }
    }

    suspend fun register(jda: JDA, logger: Logger) {
        coroutineScope {
            launch {
                for (cmd in commands) {
                    if (guildId != null) {
                        val guild = jda.getGuildById(guildId)
                        if (guild == null) {
                            throw Exception("")
                            return@launch
                        }

                        jda.getGuildById(guildId)?.upsertCommand(cmd.data)?.queue()
                        logger.info("Loaded global command: /${cmd.data.name}")
                        continue
                    }

                    jda.upsertCommand(cmd.data).queue()
                    logger.info("Loaded global command: /${cmd.data.name}")
                }
            }
        }
    }
}