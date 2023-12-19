package net.projecttl.kuma.engine.build

import net.dv8tion.jda.api.EmbedBuilder
import net.dv8tion.jda.api.entities.MessageEmbed

class NewEmbedBuilder {
    data class FieldData(
        var name: String = "",
        var value: String = "",
        var inline: Boolean = false
    )

    data class FooterData(
        var text: String? = null,
        var iconUrl: String? = null
    )

    data class Thumbnail(var url: String = "")

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
    /** Embed Thumbnail */

    fun field(field: FieldData.() -> Unit) {
        val f = FieldData().apply(field)
        embed.addField(f.name, f.value, f.inline)
    }

    fun footer(footer: FooterData.() -> Unit) {
        val f = FooterData().apply(footer)
        embed.setFooter(f.text, f.iconUrl)
    }

    fun thumbnail(thumb: Thumbnail.() -> Unit) {
        val th = Thumbnail().apply(thumb)
        embed.setThumbnail(th.url)
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
