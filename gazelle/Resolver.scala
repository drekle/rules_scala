package io.bazel.rules_scala.gazelle

import java.nio.file.{Files, Paths}

import scala.io.StdIn
import java.util.concurrent.ConcurrentLinkedQueue

import io.circe._, io.circe.generic.auto._, io.circe.parser._, io.circe.syntax._
import scala.meta.inputs.Input, scala.meta._

// const val BUFFER_SIZE = 10;

// repo_root=/Users/dlemon/repos/rules_scala;rel_package_path=gazelle;filenames=Resolver.scala
class InputObject{
    // repo_root
    var repoRoot: String = ""
    // rel_package_path
    var relPackagePath: String = ""
    // filenames
    var filenames: Array[String] = Array.empty[String]
}

class OutputObject{
    // Line output example: package=my.com.pacakge.name;main_methods=my.com.package.name.Resolver;imports=io.circe._,java.util.concurrent
    // package=my.com.pacakge.name
    var packageName: String = ""
    // main_methods=my.com.package.name.Resolver
    var mainMethods: Array[String] = Array.empty[String]
    // imports=io.circe._,java.util.concurrent
    var imports: Array[String] = Array.empty[String]
}

class InputParser(inputQueue : ConcurrentLinkedQueue[InputObject]) extends Runnable {

    def run {
        while(true) {
            val line = scala.io.StdIn.readLine()
            val kvs = line.split(";")
            val input = new InputObject()
            kvs.foreach(kv => {
                val elements = kv.split("=")
                if (elements.length == 2) {
                    elements(0) match {
                        case "repo_root" => input.repoRoot = elements(1)
                        case "rel_package_path" => input.relPackagePath = elements(1)
                        case "filenames" => {
                            var fileList = elements(1).split(",")
                            input.filenames = fileList
                        }
                    }
                }
            })
            println("handled line: " + line)
            println("filenames: " + input.filenames.length)
            inputQueue.add(input)
        }
    }
}

class Runner(workQueue : ConcurrentLinkedQueue[InputObject], outputQueue : ConcurrentLinkedQueue[OutputObject]) extends Runnable {

    def run {
        while(true) {
            val input = workQueue.poll()
            if ( null != input) {

                input.filenames.foreach(filename => {
                    println("parsing " + filename)
                    val path = Paths.get(input.repoRoot, input.relPackagePath, filename)
                    val bytes = Files.readAllBytes(path)
                    val text = new String(bytes, "UTF-8")
                    val vf = Input.VirtualFile(path.toString, text)
                    val exampleTree = vf.parse[Source].get.structure
                    println(exampleTree)
                })
            }
        }
    }

}

object Resolver {
    def main(args: Array[String]) = {
        
        val inputQueue = new ConcurrentLinkedQueue[InputObject]()
        val outputQueue = new ConcurrentLinkedQueue[OutputObject]()
        var worker = new Thread(new Runner(inputQueue, outputQueue))
        worker.start()
        var inputReader = new Thread(new InputParser(inputQueue))
        inputReader.start()

        while(true) {
            outputQueue.forEach((e) => {
                println(e.packageName)
                outputQueue.remove(e)
            })
        }
    }
}