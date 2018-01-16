## skelton class

A class used to describe project skelton information.

{
    templateDir: "", //template directory location
    ignore: [], // ignore some files or directories in templateDir
    commands: {
        before, //before creating project by skeleton description, this command would be executed.
        after, //after creating project by skeleton description, this command would be executed.
    },
    subs: [{
       skeltonPath: "",
       targetPath: "",
       contextPath: "" // like "a.b", used to pass current context's sub context into sub skeleton.
    }]
}

path

1. absolute path 
2. relative path to current file or process cwd
3. remote url

## skelton config

A config to refer skelton class, used to initialize a skeleton.

{
    skeltonPath: "",
    targetPath: "",
    context: {}
}
