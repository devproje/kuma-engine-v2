package net.projecttl.bot.test

import net.projecttl.kuma.engine.EnvConfig

object Config : EnvConfig() {
    val owner by useConfig<String>()
    val guild by useConfig<String>()
}