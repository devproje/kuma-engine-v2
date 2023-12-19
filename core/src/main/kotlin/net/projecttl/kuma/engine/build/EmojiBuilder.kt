package net.projecttl.kuma.engine.build

class EmojiBuilder(val name: String) {
    var animate = false
    var id: String = ""

    fun build(): String {
        if (id == "") {
            return ":$name:"
        }

        return "<${if (animate) "a" else ""}:$name:$id>"
    }
}