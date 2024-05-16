![icon](ico/icon.ico)

# StaticDeployment

通过字符替换将一个网页生成为多个版本。

适用于做好一张网页后，需要生成多个版本的情况。例如多语言或者只有少部分变化，又不想使用动态网页的情况。本程序用于通过配置文件托管这一流程，并且能够处理这之间的编译和转换等流程，保持原有项目的编译脚本可用并自动执行他们。

运行基于配置文件，支持 JSON 和 YAML 两种格式。操作均在脚本中完成。支持在替换前和替换后执行自定义命令，内置文件操作处理器和一些扩展功能，配置文件命令参数支持变量。

[测试](#测试) 用配置文件演示了 [操作过程](#测试脚本进行的操作) 。

## 注意事项

### 路径和大小写

- 为了防止意外的行为，本文档中所有的命令均区分大小写，请保持每个命令的大小写正确。
- Windows 中运行本程序时，带上 `.exe` 扩展名，例如 `StaticDeployment.exe` 而不是 `StaticDeployment` 。

### 临时修改文件后因出错而没有还原文件的恢复

1. 在任何对源文件的操作前，请在脚本中 `prerun` 部分中使用 `$BAK` 功能备份，并在操作结束后在脚本中 `run` 部分用 `$RES` 还原。
2. 处理脚本时如果因为错误而中断，程序会自动进入恢复模式，扫描所有经过 `$BAK` 备份而没有来得及用 `$RES` 还原的问题。如果没有自动进入恢复模式，请立即运行 `./StaticDeployment -r` 恢复您的文件，而**不要在恢复前再次尝试运行脚本**防止备份文件被覆盖。

### 配置文件安全性警告

错误的配置文件或潜在的 BUG 可能会 **损毁您的源代码！** 请在调试配置文件之前， **将您的完整的源代码备份到其他位置中！**

由于配置文件操作自由度较高（ **可以进行任意文件操作和命令操作** ），请在使用别人提供的配置文件之前 **仔细检查行为** ，在自己编写配置文件时 **仔细检查文件路径** ，以免 **电脑数据损毁！**

在本程序运行配置文件时， **请勿中断！** （包括按 `Ctrl+C` 中断），这可能导致不可意料的后果（例如您的源代码被临时修改而没有来得及还原，或者文件操作被中断而导致文件损毁）！如果中断，参考上一节。

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
2. 在自己的项目中创建本程序所使用的配置文件，配置文件写法按下面的介绍，可以是 `yaml` 也可以是 `json` 。以下均用 YAML 进行演示。
3. 运行命令: `StaticDeployment <配置文件路径>` 。

如果将 `<配置文件路径>` 省略，将自动寻找启动文件夹下的 `StaticDeployment.yaml` 作为配置文件启动。

### 回滚模式

如果操作图中发生中断，可以用下面的命令扫描所有的备份文件并恢复。

运行命令: `StaticDeployment -r <要扫描的文件夹>`

`<要扫描的文件夹>` 是可选的，默认为当前文件夹。

注意：恢复将**覆盖现有文件**，这是不可撤销的。

## 配置文件

### 配置文件示例

见： `testconfig.yaml` 。里面有详细的注释。

### 配置文件格式

#### 配置文件的基本结构

```yaml
# root: 解决方案 `Solution`
- name: 解决方案名称 1
  projects: 项目列表
    - name: 项目名称 1.1
      source: 要处理的源文件
      replaces: 替换列表
        - name: 任务名称 1.1.1
          items: 替换项
            - old: 旧字符串
              new: 新字符串
              num: 替换几次（默认 1 ）
        - name: 任务名称 1.1.2
          items: 替换项
            - old: 旧字符串
              new: 新字符串
              num: 替换几次
    - name: 项目名称 1.2
      source: 要处理的源文件
      replaces: 替换列表
        - name: 任务名称 1.2.1
          items: 替换项
            - old: 旧字符串
              new: 新字符串
              num: 替换几次
- name: 解决方案名称 2
  projects: 项目列表
    - name: 项目名称 2.1
      source: 要处理的源文件
      replaces: 替换列表
        - name: 任务名称 2.1.1
          items: 替换项
            - old: 旧字符串
              new: 新字符串
              num: 替换几次
```

#### 解决方案、项目、替换任务 都可以指定命令

- `prerun`: 操作之前运行什么
- `run`: 操作之后运行什么

格式为：

```yaml
# root: 解决方案 `Solution`
- name: 解决方案名称
  prerun: 处理解决方案之前运行什么
  projects: #项目列表
    - name: 项目名称
      prerun: 处理项目之前运行什么
      source: 要处理的源文件
      replaces: #替换列表
        - name: 任务名称
          prerun: 处理任务之前运行什么
          replace: #替换列表
            - old: 旧字符串
              new: 新字符串
              num: 替换几次
          prerun: 处理任务之后运行什么
      run: 处理项目之后运行什么
  run: 处理解决方案之后运行什么
```

#### `run` 和 `prerun` 命令的基本结构

```yaml
# root: 解决方案 `Solution`
- name: 解决方案名称
  prerun: 处理解决方案之前运行什么
    default:
      - - 命令 1
        - 命令 1 的参数 1
        - 命令 1 的参数 2
      - - 命令 2
        - 命令 2 的参数 1
        - 命令 2 的参数 2
  run: 处理解决方案之后运行什么
    default:
      - - 命令 1
```

示例:

```yaml
- name: index
  prerun:
    default:
      - - node_modules/.bin/eslint
        - --ext
        - .js,.jsx,.ts,.tsx
        - src
        - --ignore-pattern
        - "*.d.ts"
      - - node_modules/.bin/webpack
        - --mode
        - production
```

#### 自定义变量

- 可以在 `run` / `prerun` 中通过 `$SET` / `$CMDSET` / `$UNSET` 设置自定义变量；
- 然后在 `run` / `prerun` (作为外部命令参数) `items/new` (作为要替换的文本) 中通过 `$变量名` 使用自定义变量。

```yaml
# root: 解决方案 `Solution`
- name: 解决方案名称
  prerun: 处理解决方案之前运行什么
    default:
      - - $SET
        - 要创建的变量名
        - 变量内容
      - - $CMDSET
        - 要创建的变量名
        - 外部命令
        - 外部命令的参数...
  run:
    default:
      - - 外部命令
        - 外部命令
        - 可以使用$变量名
      - - $UNSET
        - 要删除的变量名
```

示例:

```yaml
- name: test1-2
  prerun:
    default:
      - - $SET
        - VAR1
        - Hello!
      - - $CMDSET
        - VAR2
        - ls
        - -ahl
  run:
    default:
      - - echo
        - $VAR1
      - - echo
        - $VAR2
      - - $UNSET
        - VAR1
      - - $UNSET
        - VAR2
```

#### 替换字符串支持变量

- 如果之前通过 [自定义变量](#自定义变量) 定义了变量，会自动将 新字符串 `new` 中的 `$变量名` 的部分替换为该变量的值。
- 可以在新字符串 `new` 中定义 `$IMPORT=` + 一个或多个文件路径（使用 `,` 分隔），将会读取这些文件的内容赋值给新字符串 `new` 。

```yaml
# root: 解决方案 `Solution`
- name: 解决方案名称
  projects: 项目列表
    - name: 项目名称
      prerun: 处理项目之前运行什么
        default:
          - - $SET
            - 变量名
            - 变量内容
      source: 要处理的源文件
      replaces: 替换列表
        - name: 任务名称
          items: 替换项
            - old: 旧字符串1
              new: 新字符串是$变量名
              num: 替换几次（默认 1 ）
            - old: 旧字符串2
              new: $IMPORT=src/2.html,src/3.html
              num: 1
```

以上配置中：

- `旧字符串1` 将被替换为 `新字符串是变量内容` 。
- `旧字符串2` 将被替换为 `src/2.html` + `src/3.html` 的文件内容。

**以下是配置文件中的各种可选项：**

### run / prerun: 系统平台可选项

- `default`: 默认的命令。如果没有指定以下特定平台，则执行这里面的命令。
- `windows`
- `linux`
- `darwin`: macOS

### run / prerun: 内置命令

`$` 开头的内置命令全都需要大写。

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
- `$MD` : 新建文件夹。
  - 参数 1 : 新的文件夹路径。支持多级创建，权限使用最近一个有效文件夹的权限。
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
- `$SET` : 设置自定义变量。可以在替换任务中的新字符串 `new` 中通过 `$变量名` 使用，也可以在后续的命令参数中使用。
  - 参数 1 : 变量名（必须，请设置长一些以免冲突，不要加 `$` ）
  - 参数 2 : 变量值（必须）
  - 示例: `["$SET", "key", "val"]`
- `$CMDSET` : 执行命令，并将命令的结果保存到变量（而不是直接输出）
  - 参数 1 : 变量名（必须，请设置长一些以免冲突，不要加 `$` ）
  - 参数 2 : 命令（必须）
  - 参数 N ... : 命令的参数
  - 示例: `["$CMDSET", "FILELIST", "ls", "-ahl"]`
- `$UNSET` : 删除自定义变量。
  - 参数 1 : 变量名（必须）
  - 示例: `["$UNSET", "key"]`
- `$RUNSLN` : 运行另一个配置文件（解决方案）
  - 参数 1 : 配置文件路径（必须）
  - 示例: `["$RUNSLN", "subsln/StaticDeployment.yaml"]`
  - 注意：将启动一个本程序的新实例，启动目录会临时修改到配置文件所在目录，然后运行配置文件。目标配置文件中的相对路径应以该配置文件所在文件夹为起始。上面的示例起始目录为 `subsln/` 。

### run / prerun: 插件命令

当 `$` 开头的命令不满足以上任意一个内置命令时，视为插件命令。支持大小写并在部分文件系统上区分大小写。

插件命令会在命令前加上本程序名，然后尝试寻找相同文件夹下的可执行文件来运行。例如 `$ZhCodeConv` 默认情况下实际会执行 `./StaticDeployment_ZhCodeConv` 。

以下为插件提供的命令，因此需要把 `./StaticDeployment_插件名` 和本程序 ( `./StaticDeployment` ) 放在同一个文件夹里面。

注意：本程序执行文件名被修改，插件的前缀名也应该同时进行修改。插件可以自由扩展，可以查看 [插件](#插件) 节了解插件的开发方法。

- `$ZhCodeConv`: 非 ASCII 变量和函数名转换。
  - 功能介绍:
    - 将所有非 ASCII 字符（例如中文）转换为指定格式。
    - 支持不处理引号之间的内容。
    - 连续的 ASCII 字符视为一个单词，如果转换后单词最前面的字符是数字，则在前面加上 `g` 。
  - 参数:
    - 参数 1 : 输入文件路径。（可选，默认为当前处理的文件 `source` ）
    - 参数 2 : 将单词转换为什么（可选）:
      - 可选项:
        - 非加密哈希函数
          - `hex`: 转换为十六进制字符串（默认值，位数: `字符数 × 6` ）
          - `crc32`: CRC-32 (IEEE) 哈希值（位数: 8 ）
          - `crc64`: CRC-64 (ISO) 哈希值（位数: 16 ）
          - `fnv1a`: FNV-1a 哈希值（位数: 8 ）
          - `adler32`: Adler-32 哈希值（位数: 8 ）
        - 加密哈希函数
          - `md5`: MD5 哈希值（位数: 32 ）
          - `sha1`: MD5 哈希值（位数: 40 ）
          - `sha256`: MD5 哈希值（位数: 64 ）
          - `sha512`: MD5 哈希值（位数: 128 ）
      - 如果首位不是字母则 `位数 + 1`
      - 支持指定迭代次数，在后面加 `*数字` 即可。例如 `md5*32` 为将字符串用 `md5` 计算 32 次。
    - 参数 3 : 输出文件路径。（可选，默认为直接输出，可使用 `$CMDSET` 收取成变量）
    - 参数 4 : 这些符号之间的内容不处理。（可选）
      - 默认为 `"'` ，即 `" - "` 和 `' - '` 之间的内容不做处理。
      - 仅作外层首个遇到的符号匹配，例如 `aa"bb'cc'"dd` 中 `bb` 和 `cc` 都会被处理。
      - 留空则禁用此功能，所有的都会被处理。在代码中的字符串皆不包含中文且需要运行在字符串中的代码时，此处建议留空。
  - 示例:
    - 源文件 `test1.js` : `var 中文变量 = '测试'; function 中文函数() { 中文变量 = "Hello, 世界!"; console.log(中文变量); } 中文函数();`
    - 命令: `["$ZhCodeConv", "test1.js", "md5", "dist/test1.js"]`
    - 输出文件 `dist/test1.js` : `var e4b8ade69687e58f98e9878f = '测试'; function e4b8ade69687e587bde695b0() { e4b8ade69687e58f98e9878f = "Hello, 世界!"; console.log(e4b8ade69687e58f98e9878f); } e4b8ade69687e587bde695b0();`
- `$Minify`: 使用 [Minify](https://github.com/tdewolff/minify) 库压缩代码。
  - 参数:
    - 参数 1 : 输入文件路径。
    - 参数 2 : 输出文件路径。（可选，默认为直接输出，可使用 `$CMDSET` 收取成变量）
    - 参数 3 : 要压缩的代码文件类型（可选，根据文件扩展名自动判断）
      - 支持的选项: `html`, `css`, `js`, `json`, `svg`, `xml`
  - 示例:
    - 源文件 `test1.html` : `<div></div>\n\n<div></div>`
    - 命令: `["$Minify", "test1.html", "dist/test1.html", "html"]`
    - 输出文件 `dist/test1.html` : `<div></div><div></div>`
- `$Join`: 将多个文件追加合并到一个文件中。
  - 参数:
    - 参数 1 : 要保存合并后内容的文件。如果文件已经有内容则追加，如果文件不存在则新建。
    - 参数 ... : 要合并进 参数1 文件中的文件，每个文件路径一个文件。
  - 示例:
    - 源文件 `1.txt`:`111`;  `2.txt`:`222`;  `3.txt`:`333`
    - 命令: `["$Join", "1.txt", "2.txt", "3.txt"]`
    - 输出文件 `1.txt` : `111222333`

### run / prerun: 变量列表

`$` 开头的内置变量列表全都需要大写。

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

根目录为主程序的源码，每个大写字母开头的文件夹是插件的源码。

1. 确保已经安装了 GO 环境，并且不低于 `go.mod` 中的版本。
2. 确保系统中已经安装了以下命令: `go`, `go generate`, `go build`, `openssl`, `7z` 。
3. 在当前文件夹运行 `go get` 。
4. 运行全平台编译脚本:

- Windows 系统 (cmd): 在当前文件夹运行 `build.bat` 。
- 非 Windows 系统 (bash): 在当前文件夹运行 `chmod +x build.sh && ./build.sh` 。
- 编译后的文件将以 `程序名_系统_平台.7z` 的格式存放在主程序和每个插件的 `./bin` 文件夹中。
  - 还会以 `程序名` 的格式在 `GOPATH/bin` 中存放一份本地平台版本以便使用。

5. 编译完成后主编译脚本会暂停，等待其他插件编译进程结束后可以手动继续。
6. 可以将每个插件文件夹中的 `bin` 文件夹合并到主程序的 `bin` 文件夹。

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

## 插件

### 运行周期

以 Windows 为例，当用户在 主程序(StaticDeployment) 脚本中调用 `$命令名` 时，主程序会执行 `StaticDeployment.exe` 相同文件夹下的 `StaticDeployment_本插件编译后的文件名.exe` （如果主程序被重命名，则以新名称为准）。

例如，用户在主程序脚本中调用 `$PluginDemo` 时，主程序会执行 `StaticDeployment_PluginDemo.exe` ，即本插件最终输出的执行文件名，并放在和主程序执行文件同一个文件夹下。

插件程序在运行步骤的位置为下面的 `2.` 和 `3.`：

1. 主程序备份正在修改的文件
2. **主程序运行本插件程序**
3. **本插件程序通过 fmt 将结果输出给主程序 或者 直接改好主程序通过参数提供的文件路径**
4. 主程序根据插件程序改好的文件，或者插件程序提供的内容，进行编译等操作
5. 主程序恢复备份的文件

因此对用户的源文件进行覆盖修改是安全的也是通常做法。

### 向主程序回复信息

- 如果需要输出日志信息，使用 `log` 包，内容会直接进行输出；
- 如果需要输出返回值，使用 `fmt` 包，内容会被主程序收取（如果用户在脚本中要求）。

```go
log.Println(os.Args[0], "测试插件开始运行……")
```

如果本插件执行失败，退出时需返回 > 0 的值。主程序遇到非 0 的返回值会停止执行脚本。

```go
if len(os.Args) <= 1 {
  log.Println(os.Args[0], "错误: 未指定参数。")
  os.Exit(1)
  return
}
```

### 从主程序接收参数

用户可以在主程序脚本中传递参数，例如 `$PluginDemo 参数1 参数2 参数...`，本插件会接收到 `参数0(本插件名) 参数1 参数2 参数...`。

```go
var key string = os.Args[1]
```

例如上面这行代码，用户在脚本中提供 `$PluginDemo key1` 时, `var key` 的值即为 `key1` 。

### 示例代码

假设这里有一些数据，根据用户的请求 `var key = os.Args[1]` 拿出 `value` ：

```go
var data = map[string]string{
  "key1": "value1",
  "key2": "value2",
}
```

假设这里有个 HTML 代码片段，主程序需要拿这个 HTML 代码片段去合成整个 HTML 。

var html =

```html
<div>这是一个插件示例，根据 %KEY% ，我们找到了: %DATA%</div>
```

本插件将取到的数据填充到 HTML 代码片段中:

```go
html = strings.ReplaceAll(html, "%KEY%", key)
if value, ok := data[key]; ok {
  html = strings.ReplaceAll(html, "%DATA%", value)
} else {
  html = strings.ReplaceAll(html, "%DATA%", "什么都没找到！")
}
```

将 HTML 代码片段返回给主程序。主程序会如实接收 fmt.Println 的内容，进行进一步操作。

如果主程序脚本中没有设置接收数据，则 fmt.Println 会直接被输出（例如本插件是直接操作文件的程序）。

```go
fmt.Println(html)
```

在上面的程序中，用户在脚本中提供 `$PluginDemo key1` ，可以得到 `<div>这是一个插件示例，根据 key1 ，我们找到了: value1</div>` 的结果，并存储在一个自定义变量中。

#### 由插件程序直接读写文件

可以由用户自行提供一个文件路径来由本插件打开或存储数据：

如果插件的功能需要读取一个文件，可以将文件路径作为参数传递给本插件，本插件会读取文件内容并返回给主程序。

```go
if len(os.Args) > 2 {
  // 根据用户提供的路径读取文件
  var path = os.Args[2]
  data, err := os.ReadFile(path)
  if err != nil {
    log.Println(os.Args[0], "错误: 无法读取文件", path, err)
    os.Exit(1)
    return
  }
  // 将读取到的数据填充到 HTML 代码片段中
  html = strings.ReplaceAll(html, "%DATA%", string(data))
}
```

如果插件的功能需要将数据写入到一个文件，可以将文件路径作为参数传递给本插件，本插件会将数据写入到文件中。

```go
if len(os.Args) > 3 {
  var path = os.Args[3]
  // 例如，将 HTML 代码片段写入到文件中
  err := os.WriteFile(path, []byte(html), 0644)
  if err != nil {
    log.Println(os.Args[0], "错误: 无法写入文件", path, err)
    os.Exit(1)
    return
  }
} else {
  // 可以判断用户是否给与另存为参数，如果没有，再通过输出来返回数据。
  fmt.Println(html)
}
// 本插件执行完毕，退出时需返回 0 的值。告诉主程序可以继续运行。
os.Exit(0)
```

在上面的程序中，用户在脚本中提供 `$PluginDemo key2 1.html 1.html` ，并且 `1.html` 的内容为 `<div>这是一个插件示例，根据 %KEY% ，我们找到了: %DATA%</div>` , 可以得到 `<div>这是一个插件示例，根据 key2 ，我们找到了: value2</div>` 的结果，并且会将这个结果覆盖写入到 1.html 文件中。

### 完整插件示例

插件的基本示例见 `PluginDemo` 中的 Go 程序，已包含注释说明。

## 协议

Copyright (c) 2024 KagurazakaYashi StaticDeployment is licensed under Mulan PSL v2. You can use this software according to the terms and conditions of the Mulan PSL v2. You may obtain a copy of Mulan PSL v2 at: http://license.coscl.org.cn/MulanPSL2 THIS SOFTWARE IS PROVIDED ON AN “AS IS” BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE. See the Mulan PSL v2 for more details.
