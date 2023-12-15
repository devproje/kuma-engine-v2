package net.projecttl.kuma.engine.`object`

import net.dv8tion.jda.api.Permission
import net.dv8tion.jda.api.interactions.commands.DefaultMemberPermissions
import net.dv8tion.jda.api.interactions.commands.OptionType
import net.dv8tion.jda.api.interactions.commands.build.CommandData
import net.dv8tion.jda.api.interactions.commands.build.OptionData
import net.dv8tion.jda.internal.interactions.CommandDataImpl

class CommandDataBuilder {
    data object NewOption {
        var name: String = ""
        var description: String = ""
        var type: OptionType = OptionType.STRING
        var choices: List<NewChoice> = listOf()
        var required: Boolean = false
    }

    data object NewChoice {
        var name: String = ""
        var value: String = ""
    }

    var name: String = ""
    var description: String = ""
    private val perms = mutableListOf<Permission>()
    private val options = mutableListOf<OptionData>()

    fun option(rawOpt: NewOption.() -> Unit) {
        val opt = NewOption.apply(rawOpt)
        val data = OptionData(NewOption.type, NewOption.name, NewOption.description, NewOption.required)
        if (NewOption.choices.isNotEmpty()) {
            NewOption.choices.forEach { choice ->
                data.addChoice(NewChoice.name, NewChoice.value)
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
