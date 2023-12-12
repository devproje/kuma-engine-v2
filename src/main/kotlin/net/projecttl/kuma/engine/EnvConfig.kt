package net.projecttl.kuma.engine

import io.github.cdimascio.dotenv.Dotenv
import kotlin.reflect.KProperty

class EnvConfig {
    fun <T> useConfig(): EnvDelegate<T> {
        return EnvDelegate()
    }

    val token by useConfig<String>()
}

@Suppress("UNCHECKED_CAST")
class EnvDelegate<T> {
    private val env = Dotenv.load()

    operator fun getValue(thisRef: Any, property: KProperty<*>): T {
        val propertyName = property.name.uppercase()
        return env[propertyName] as T
    }
}