plugins {
    id("com.github.johnrengelman.shadow") version "8.1.1"
    application
}

group = "net.projecttl"
version = "1.0.0"

repositories {
    mavenCentral()
}

dependencies {
    implementation(project(":core"))
    testImplementation("org.jetbrains.kotlin:kotlin-test")
}

tasks.test {
    useJUnitPlatform()
}

application {
    mainClass = "net.projecttl.bot.example.MainKt"
}

kotlin {
    jvmToolchain(17)
}
