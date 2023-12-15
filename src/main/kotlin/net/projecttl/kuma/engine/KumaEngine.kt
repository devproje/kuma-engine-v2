package net.projecttl.kuma.engine

import kotlinx.coroutines.DelicateCoroutinesApi
import kotlinx.coroutines.GlobalScope
import kotlinx.coroutines.coroutineScope
import kotlinx.coroutines.launch
import net.dv8tion.jda.api.JDABuilder
import net.dv8tion.jda.api.hooks.ListenerAdapter
import net.dv8tion.jda.api.requests.GatewayIntent
import net.projecttl.kuma.engine.command.CommandHandler
import org.slf4j.LoggerFactory

class KumaEngine(token: String, indent: List<GatewayIntent> = listOf()) {
    private val builder = JDABuilder.createDefault(token, indent)
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

    @OptIn(DelicateCoroutinesApi::class)
    suspend fun build() {
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
