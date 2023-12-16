package net.projecttl.kuma.engine

import kotlinx.coroutines.DelicateCoroutinesApi
import kotlinx.coroutines.GlobalScope
import kotlinx.coroutines.coroutineScope
import kotlinx.coroutines.launch
import net.dv8tion.jda.api.JDABuilder
import net.dv8tion.jda.api.events.interaction.command.SlashCommandInteractionEvent
import net.dv8tion.jda.api.hooks.ListenerAdapter
import net.dv8tion.jda.api.requests.GatewayIntent
import net.dv8tion.jda.api.utils.cache.CacheFlag
import net.projecttl.kuma.engine.command.CommandExecutor
import net.projecttl.kuma.engine.command.CommandHandler
import net.projecttl.kuma.engine.`object`.CommandDataBuilder
import org.slf4j.LoggerFactory

class KumaEngine(token: String, indents: List<GatewayIntent> = listOf(), flags: List<CacheFlag> = listOf()) {
    private val builder = JDABuilder.createDefault(token, indents)
        .enableCache(flags)
    private val handlers = mutableListOf<ListenerAdapter>()
    private val commands = mutableListOf<CommandHandler>()

    private val logger = LoggerFactory.getLogger("KumaEngine")

    fun addCommandHandler(vararg command: CommandHandler) {
        handlers.addAll(command)
        commands.addAll(command)
    }

    fun dropCommandHandler(command: CommandHandler) {
        if (!handlers.contains(command) && !commands.contains(command)) {
            return
        }

        handlers.remove(command)
        commands.remove(command)
    }

    fun addHandler(vararg handler: ListenerAdapter) {
        handlers.addAll(handler)
    }

    fun dropHandler(handler: ListenerAdapter) {
        if (!handlers.contains(handler)) {
            return
        }

        handlers.remove(handler)
    }

    private object KumaInfo : CommandHandler("kumainfo") {
        private final val LOGO = "https://github.com/devproje/kuma-engine/raw/master/assets/kuma-engine-logo.png"
        init {
            addCommands(KumaInfoCommand)
        }

        private object KumaInfoCommand : CommandExecutor {
            override val data = CommandDataBuilder().apply {
                name = "kumainfo"
                description = "check kuma engine version"
            }

            override fun execute(event: SlashCommandInteractionEvent) {

            }
        }
    }

    @OptIn(DelicateCoroutinesApi::class)
    suspend fun build() {
        addCommandHandler(KumaInfo)

        coroutineScope {
            launch {
                val jda = builder.build()

                GlobalScope.launch {
                    commands.forEach {
                        it.register(jda, logger)
                    }
                }
            }
        }
    }
}
