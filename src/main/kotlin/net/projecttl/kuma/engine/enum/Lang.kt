package net.projecttl.kuma.engine.enum

import com.google.gson.Gson
import java.io.File

enum class Lang(val code: String) {
    KO_KR("ko_kr"),
    EN_US("en_us");

    fun load(): LangPath {
        val res = this::class.java.getResource("/lang/${code}.json")
            ?: throw NullPointerException("$code is not supported language")
        val file = File(res.toURI())

        return Gson().fromJson(file.reader(), LangPath::class.java)
    }
}

data class Command(
    val error: String
)

data class LangPath(
    val command: Command
)
