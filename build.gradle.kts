plugins {
    kotlin("jvm") version "1.9.21"
}

group = "net.projecttl"
version = "2.0.0"

val jda_version: String by project
val mongodb_version: String by project
val exposed_version: String by project
val coroutine_version: String by project

allprojects {
    apply(plugin = "org.jetbrains.kotlin.jvm")

    kotlin {
        jvmToolchain(17)
    }

    repositories {
        mavenCentral()
    }
}


subprojects {
    repositories {
        mavenCentral()
    }

    dependencies {
        implementation(kotlin("stdlib"))
        implementation("net.dv8tion:JDA:${jda_version}")
        implementation("com.google.code.gson:gson:2.10.1")
        implementation("io.github.cdimascio:dotenv-kotlin:6.4.1")
        implementation("org.apache.logging.log4j:log4j-slf4j-impl:2.19.0")
        implementation("org.mongodb:mongodb-driver-kotlin-coroutine:${mongodb_version}")
        implementation("org.jetbrains.kotlinx:kotlinx-coroutines-core:${coroutine_version}")

        testImplementation(kotlin("test"))
    }
}
