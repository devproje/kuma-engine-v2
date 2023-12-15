package net.projecttl.kuma.engine.`object`

import net.dv8tion.jda.api.EmbedBuilder
import net.dv8tion.jda.api.entities.MessageEmbed

class NewEmbedBuilder {
    data object NewFieldData {
        var name: String = ""
        var value: String = ""
        var inline: Boolean = false
    }

    data object NewFooterData {
        var text: String? = null
        var iconUrl: String? = null
    }

    private val embed = EmbedBuilder()

    /** Embed Image Size */
    var size: Int = 2048
    /** Embed Color */
    var color: Int = 0x1F
    /** Embed Title */
    var title: String? = null
    /** Embed Image*/
    var image: String? = null
    /** Embed Description */
    var description: String? = null

    fun field(field: NewFieldData.() -> Unit) {
        val f = NewFieldData.apply(field)
        embed.addField(f.name, f.value, f.inline)
    }

    fun footer(footer: NewFooterData.() -> Unit) {
        val f = NewFooterData.apply(footer)
        embed.setFooter(f.text, f.iconUrl)
    }

    fun build(): MessageEmbed {
        apply()
        return embed.build()
    }

    private fun apply() {
        embed.setTitle(title)
        embed.setDescription(description)
        if (image != null) {
            image += "?size=${size}"
        }

        embed.setImage(image)
        embed.setColor(color)
    }
}
