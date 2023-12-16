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
import net.projecttl.kuma.engine.`object`.NewEmbedBuilder
import org.slf4j.LoggerFactory
import java.io.File

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

    companion object {
        fun version(): String {
            return File(this::class.java.getResource("/version.txt")!!.toURI()).readText()
        }
    }

    private object KumaInfo : CommandHandler() {
        private const val LOGO = "https://github.com/devproje/kuma-engine/raw/master/assets/kuma-engine-logo.png"
        init {
            addCommands(KumaInfoCommand)
        }

        private object KumaInfoCommand : CommandExecutor {
            override val data = CommandDataBuilder().apply {
                name = "kumainfo"
                description = "kuma engine system information"
            }

            override fun execute(event: SlashCommandInteractionEvent) {
                val embed = NewEmbedBuilder().apply {
                    title = ":dart: **KumaInfo**"
                    description = "KumaEngine system information"
                    thumbnail {
                        url = LOGO
                    }

                    field {
                        name = ":electric_plug: **ENGINE VERSION**"
                        value = "`${version()}`"
                        inline = true
                    }

                    field {
                        name = ":page_facing_up: **KOTLIN VERSION**"
                        value = "`${KotlinVersion.CURRENT}`"
                        inline = true
                    }

                    field {
                        name = ":ping_pong: **API LATENCY**"
                        value = "${event.jda.gatewayPing}ms"
                        inline = true
                    }

                    field {
                        name = ":desktop: **OS**"
                        value = "${System.getProperty("os.name")}/${System.getProperty("os.arch")}"
                    }
                }.build()

                event.replyEmbeds(embed).queue()
//                    {
//                        Name:   fmt.Sprintf("%s **OS**", emoji.Desktop),
//                        Value:  fmt.Sprintf("`%s/%s`", runtime.GOOS, runtime.GOARCH),
//                        Inline: true,
//                    },
//                    {
//                        Name:   fmt.Sprintf("%s **BOT SERVERS**", emoji.Satellite),
//                        Value:  fmt.Sprintf("`%d`", len(session.State.Guilds)),
//                        Inline: true,
//                    },
//                    {
//                        Name:   fmt.Sprintf("%s **SYSTEM PID**", emoji.FileFolder),
//                        Value:  fmt.Sprintf("`%d`", os.Getpid()),
//                        Inline: true,
//                    },
//                },
//                Color: rand.Intn(0xFFFFFF),
//                Footer: &discordgo.MessageEmbedFooter{
//                    Text:    event.Member.User.String(),
//                    IconURL: event.Member.User.AvatarURL("512x512"),
//                },
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
