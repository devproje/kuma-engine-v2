package org.example.net.projecttl.bot.example

import net.projecttl.kuma.engine.EnvConfig

object Config : EnvConfig() {
    val owner = useConfig<String>()
}