package net.projecttl.kuma.engine

import kotlinx.coroutines.DelicateCoroutinesApi
import kotlinx.coroutines.GlobalScope
import kotlinx.coroutines.coroutineScope
import kotlinx.coroutines.launch
import net.dv8tion.jda.api.JDABuilder
import net.dv8tion.jda.api.hooks.ListenerAdapter
import net.dv8tion.jda.api.requests.GatewayIntent
import net.projecttl.kuma.engine.command.CommandHandler

class KumaEngine(token: String, indent: List<GatewayIntent> = listOf()) {
    private val builder = JDABuilder.createDefault(token, indent)
    private val handlers = mutableListOf<ListenerAdapter>()
    private val commands = mutableListOf<CommandHandler>()

    @OptIn(DelicateCoroutinesApi::class)
    suspend fun build() {
        coroutineScope {
            launch {
                val jda = builder.build()

                GlobalScope.launch {
                    commands.forEach {
                        it.register(jda)
                    }
                }
            }
        }
    }
}
