package net.projecttl.kuma.engine.`object`

import net.dv8tion.jda.api.Permission
import net.dv8tion.jda.api.entities.MessageEmbed
import net.dv8tion.jda.api.entities.MessageEmbed.Thumbnail
import net.dv8tion.jda.api.interactions.commands.DefaultMemberPermissions
import net.dv8tion.jda.api.interactions.commands.OptionType
import net.dv8tion.jda.api.interactions.commands.build.CommandData
import net.dv8tion.jda.api.interactions.commands.build.OptionData
import net.dv8tion.jda.internal.interactions.CommandDataImpl

class CommandDataBuilder {
    data class Option(
        var name: String = "",
        var description: String = "",
        var type: OptionType = OptionType.STRING,
        var choices: List<Choice> = listOf(),
        var required: Boolean = false
    )

    data class Choice(
        var name: String = "",
        var value: String = ""
    )

    var name: String = ""
    var description: String = ""
    private val perms = mutableListOf<Permission>()
    private val options = mutableListOf<OptionData>()

    fun option(rawOpt: Option.() -> Unit) {
        val opt = Option().apply(rawOpt)
        val data = OptionData(opt.type, opt.name, opt.description, opt.required)
        if (opt.choices.isNotEmpty()) {
            opt.choices.forEach { choice ->
                data.addChoice(choice.name, choice.value)
            }
        }

        options.add(data)
    }

    fun perm(permission: Permission) {
        perms += permission
    }

    fun build(): CommandData {
        val data = CommandDataImpl(name, description)
        if (options.isNotEmpty()) {
            data.addOptions(options)
        }

        if (perms.isNotEmpty()) {
            perms.forEach { perm ->
                data.setDefaultPermissions(DefaultMemberPermissions.enabledFor(perm))
            }
        }

        return CommandData.fromData(data.toData())
    }
}
