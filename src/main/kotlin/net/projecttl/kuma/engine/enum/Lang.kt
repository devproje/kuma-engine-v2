package net.projecttl.kuma.engine.enum

enum class Lang(val code: String) {
    KO_KR("ko_kr"),
    EN_US("en_us");

    fun load(lang: Lang, path: String): Any {
    }
}