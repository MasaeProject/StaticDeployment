![icon](ico/icon.ico)

# StaticDeployment

通过字符替换将一个网页生成为多个版本。

适用于做好一张网页后，需要生成多个版本的情况。例如多语言或者只有少部分变化，又不想使用动态网页的情况。本程序用于通过配置文件托管这一流程，并且能够处理这之间的编译和转换等流程，保持原有项目的编译脚本可用并自动执行他们。

运行基于配置文件，支持 JSON 和 YAML 两种格式。操作均在脚本中完成。支持在替换前和替换后执行自定义命令，内置文件操作处理器和一些扩展功能，配置文件命令参数支持变量。

[测试](#测试) 用配置文件演示了 [操作过程](#测试脚本进行的操作) 。

## 注意事项

错误的配置文件或潜在的 BUG 可能会 **损毁您的源代码！** 请在调试配置文件之前， **将您的完整的源代码备份到其他位置中！**

由于配置文件操作自由度较高（ **可以进行任意文件操作和命令操作** ），请在使用别人提供的配置文件之前 **仔细检查行为** ，在自己编写配置文件时 **仔细检查文件路径** ，以免 **电脑数据损毁！**

在本程序运行配置文件时， **请勿中断！** （包括按 `Ctrl+C` 中断），这可能导致不可意料的后果（例如您的源代码被临时修改而没有来得及还原，或者文件操作被中断而导致文件损毁）！

使用本工具造成的源代码损坏或其他损坏， **作者概不负责。** 总之请注意备份您的文件。

## 示例

例如源文件 `src/index.html`

```html
<!DOCTYPE html>
<body>
  Hello, World!
</body>
</html>
```

转换后可以输出为：

`out/zh-cn/index.html`

```html
<!DOCTYPE html>
<body>
  你好世界！
</body>
</html>
```

`out/ja-jp/index.html`

```html
<!DOCTYPE html>
<body>
  こんにちは世界！
</body>
</html>
```

以及它自身到 `out/en-us/index.html`

## 使用

1. 将本执行文件放在 `PATH` 环境变量里或拷贝到项目文件夹下以便能在别的地方输入该命令运行本程序。
2. 在自己的项目中创建本程序所使用的配置文件，配置文件写法按下面的介绍，可以是 `yaml` 也可以是 `json` 。
3. 运行命令: `StaticDeployment <配置文件路径>` 。

## 配置文件

### 配置文件示例

见： `testconfig.yaml` 。里面有详细的注释。

### 配置文件格式

配置文件的基本结构：

```yaml
# root: 解决方案 `Solution`
- name: 解决方案名称 1
  - projects: 项目列表
    - name: 项目名称 1.1
      source: 要处理的源文件
      replace: 替换列表
        - name: 任务名称 1.1.1
          replace: 替换列表
            - old: 旧字符串
              new: 新字符串
              num: 替换几次
        - name: 任务名称 1.1.2
          replace: 替换列表
            - old: 旧字符串
              new: 新字符串
              num: 替换几次
    - name: 项目名称 1.2
      source: 要处理的源文件
      replace: 替换列表
        - name: 任务名称 1.2.1
          replace: 替换列表
            - old: 旧字符串
              new: 新字符串
              num: 替换几次
- name: 解决方案名称 2
  - projects: 项目列表
    - name: 项目名称 2.1
      source: 要处理的源文件
      replace: 替换列表
        - name: 任务名称 2.1.1
          replace: 替换列表
            - old: 旧字符串
              new: 新字符串
              num: 替换几次
```

解决方案、项目、替换任务 都可以指定：

- `prerun`: 操作之前运行什么
- `run`: 操作之后运行什么

格式为：

```yaml
# root: 解决方案 `Solution`
- name: 解决方案名称
  - prerun: 处理解决方案之前运行什么
  - projects: 项目列表
    - prerun: 处理项目之前运行什么
    - name: 项目名称
      source: 要处理的源文件
      replace: 替换列表
        - prerun: 处理任务之前运行什么
        - name: 任务名称
          replace: 替换列表
            - old: 旧字符串
              new: 新字符串
              num: 替换几次
        - prerun: 处理任务之后运行什么
    - run: 处理项目之后运行什么
  - run: 处理解决方案之后运行什么
```

**以下是配置文件中的各种可选项：**

### run / prerun: 系统平台可选项

- `default`: 默认的命令。如果没有指定以下特定平台，则执行这里面的命令。
- `windows`
- `linux`
- `darwin`: macOS

### run / prerun: 内置命令

- `$CMDDIR`: 本段命令组中，接下来要运行的所有**外部**命令，都会基于这个文件夹。参数如下：
  - 参数 1 : 外部命令的工作目录。（可选，默认为当前本程序运行目录）
  - 示例: `["$CMDDIR", "outdir/project1"]`
- `$BAK` : 将当前处理的文件备份为 `文件名 + .StaticDeploymentBackup` 以便对当前文件操作后执行正常的项目编译脚本，应与 `$RES` 成对出现。
  - 参数 1 : 要备份的文件路径。（可选，默认为当前处理的文件 `source` ）
  - 示例: `["$BAK", "test1.js"]`
- `$RES` : 还原使用 `$BAK` 备份的文件。
  - 参数 1 : 解决方案名，将还原此解决方案中的全部文件，除非指定了下面的更细化参数。（可选，默认为还原此解决方案中的此项目中的此任务中的全部文件）
  - 参数 2 : 项目名，将还原此解决方案中的此项目中的全部文件，除非指定了下面的更细化参数。（可选）
  - 参数 3 : 任务名，将还原此解决方案中的此项目中的此任务中的全部文件。（可选）
  - 示例: `["$RES", "Solution1", "Project1", "ReplaceJob1"]`
- `$CP` : 复制文件或文件夹。
  - 参数 1 : 路径:
    - 如果提供参数 2 ，则此参数为将当前处理的文件操作到哪里（目标文件路径）；
    - 如果不提供参数 2 ，则此参数为源文件路径。
  - 参数 2 : 目标文件路径。（可选）
  - 示例: `["$CP", "test1.js", "dist/test1.js"]`
- `$MV` : 移动文件或文件夹（直接移动）。
  - 参数和示例同 `$CP` 。
- `$SMV` : 移动文件或文件夹（先复制文件，再删除源文件）。
  - 参数和示例同 `$CP` 。
- `$REN` : 重命名文件或文件夹。
  - 参数和示例同 `$CP` 。
- `$RM` : 删除文件或文件夹（直接删除）。
  - 参数 1 : 路径（必须）
  - 示例: `["$RM", "dist/test1.js"]`
- `$SRM` : 删除文件或文件夹（先覆盖填充文件，再删除）。
  - 参数和示例同 `$RM` 。

#### 扩展功能

- `$ZHCODECONV`: 非 ASCII 变量和函数名转换。
  - 功能介绍:
    - 将所有非 ASCII 字符（例如中文）转换为指定格式。
    - 不处理引号之间的内容。
    - 连续的 ASCII 字符视为一个单词，如果转换后单词最前面的字符是数字，则在前面加上 `g` 。
  - 参数:
    - 参数 1 : 输入文件路径。（可选，默认为当前处理的文件 `source` ）
    - 参数 2 : 将单词转换为什么（可选）:
      - 可选项:
        - `hex`: 转换为十六进制字符串（默认值，位数: `字符数 × 6` ）
        - `md5`: MD5 哈希值（位数: 32 ）
        - `sha1`: MD5 哈希值（位数: 40 ）
        - `sha256`: MD5 哈希值（位数: 64 ）
        - `sha512`: MD5 哈希值（位数: 128 ）
      - 如果首位不是字母则 `位数 + 1`
      - 支持指定迭代次数，在后面加 `*数字` 即可。例如 `md5*32` 为将字符串用 `md5` 计算 32 次。
    - 参数 3 : 输出文件路径。（可选，默认为当前处理的文件 `source` ，即**直接覆盖**，覆盖操作需要配合 `$BAK` 和 `$RES` 使用）
  - 示例:
    - 源代码: `var 中文变量 = '测试'; function 中文函数() { 中文变量 = "Hello, 世界!"; console.log(中文变量); } 中文函数();`
    - 命令: `["$ZHCODECONV", "test1.js", "md5", "dist/test1.js"]`
    - 输出: `var e4b8ade69687e58f98e9878f = '测试'; function e4b8ade69687e587bde695b0() { e4b8ade69687e58f98e9878f = "Hello, 世界!"; console.log(e4b8ade69687e58f98e9878f); } e4b8ade69687e587bde695b0();`
- `$MINIFY`: 使用 [Minify](https://github.com/tdewolff/minify) 库压缩代码。
  - 参数 1 : 输入文件路径。（可选，默认为当前处理的文件 `source` ）
  - 参数 2 : 要压缩的代码文件类型（可选，根据文件扩展名自动判断）
    - 支持的选项: `html`, `css`, `js`, `json`, `svg`, `xml`
  - 参数 3 : 输出文件路径。（可选，默认为当前处理的文件 `source` ，即**直接覆盖**，覆盖操作需要配合 `$BAK` 和 `$RES` 使用）

### run / prerun: 变量列表

- `$SOLUTION`: 当前解决方案名称
  - 例如: `Solution1`
- `$PROJECT` : 当前项目名称
  - 例如: `Project1`
- `$JOBNAME` : 当前替换任务名称
  - 例如: `ReplaceJob1`
- `$SRC` : 当前处理的源文件 `source` 的文件路径
  - 例如: `codes/project1/src/test1.js`
- `$SRCFILE` : 当前处理的源文件 `source` 的文件名
  - 例如: `test1.js`
- `$SRCNAME` :当前处理的源文件 `source` 的文件名（不带扩展名）
  - 例如: `test1`
- `$SRCEXT` : 当前处理的源文件 `source` 的扩展名
  - 例如: `js`
- `$SRCDIR` : 当前处理的源文件 `source` 的文件夹路径
  - 例如: `codes/project1/src`
- `$SRCDIRNAME` : 当前处理的源文件 `source` 的当前文件夹名
  - 例如: `src`
- 示例:
  - `out/$SRCDIRNAME/$JOBNAME/$PROJECT/$SRCNAME.$SRCEXT`

## 编译

1. 确保已经安装了 GO 环境，并且不低于 `go.mod` 中的版本。
2. 确保系统中已经安装了以下命令: `go`, `go generate`, `go build`, `openssl`, `7z` 。
3. 在当前文件夹运行 `go get` 。
4. 运行全平台编译脚本:

- Windows 系统 (cmd): 在当前文件夹运行 `build.bat` 。
- 非 Windows 系统 (bash): 在当前文件夹运行 `chmod +x build.sh && ./build.sh` 。

编译后的文件将以 `程序名_系统_平台.7z` 的格式存放在 `./bin` 文件夹中。还会以 `程序名` 的格式在 `GOPATH/bin` 中存放一份本地平台版本以便使用。

## 测试

在 `test/data` 中有个三个测试用前端项目，其中 `testconfig.yaml` 或 `testconfig.json` 中包含有完整的测试流程。

1. 请确保系统中已经安装 `node` 并且 `npm` 命令可用: 测试项目用到了 [Node.js](https://nodejs.org/) 。
2. 下载引用的子项目，在本仓库根目录运行 `git submodule init` ，以后如果更新本仓库后，需要运行 `git submodule update --remote` 。
3. 前往 `cd test/data/test-project3` 安装这个测试项目的依赖 `npm install` 并测试编译 `npm run release` 。
4. 回到仓库根目录 `cd ../../../` 完成上述 [编译](#编译) 步骤。
5. 在本仓库根目录运行命令: `go run . test/testconfig.yaml` 。

输出的文件在 `./out` 文件夹中。

### 测试脚本进行的操作

#### ./test/data/test-project1

该测试进行了一个最基本的 HTML 替换并导出的流程：

1. 备份源代码文件
2. 进行修改代码和替换文本任务
3. 复制修改后文件到输出文件夹
4. 还原备份源代码文件
5. 开始处理下一个源代码文件（或者下一个副本）

#### ./test/data/test-project2

该测试是 TypeScript 写的代码，因此每个源文件需要分别编译才能使用：

1. 备份源代码文件
2. 进行修改代码和替换文本任务
3. 执行目标项目的编译脚本
4. 复制编译后文件到输出文件夹
5. 还原备份源代码文件
6. 开始处理下一个源代码文件（或者下一个副本）

#### ./test/data/test-project3

该测试是一个完整 webpack 项目，所有文件必须一起编译：

1. 执行目标项目的编译脚本
2. 备份编译后代码文件
3. 进行修改代码和替换文本任务
4. 复制修改后文件到输出文件夹
5. 还原编译后代码文件
6. 开始处理下一个编译后代码文件（或者下一个副本）

## 协议

Copyright (c) 2024 KagurazakaYashi StaticDeployment is licensed under Mulan PSL v2. You can use this software according to the terms and conditions of the Mulan PSL v2. You may obtain a copy of Mulan PSL v2 at: http://license.coscl.org.cn/MulanPSL2 THIS SOFTWARE IS PROVIDED ON AN “AS IS” BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE. See the Mulan PSL v2 for more details.
