

# 一、目标与背景

通过系统运行过程产生的日志、遥测数据或某些业务数据的分析实现对系统当前运行状态、历史状态的变迁的体现。

整个流程包括数据生成、数据输出、数据采集、数据入库前解析处理、数据入库、数据搜索、数据搜索时解析、数据可视化和数据分析等。

本文主要注重于数据输出部分流程。数据输出模块对接数据生成和数据采集模块，需要考虑上下两个模块的特性和需求。

![](./images/span-log%20life.png)



# 二、数据输出相邻模块需求分析

## 2.1 生成模块数据特性

数据生成模块主要会生成metric、trace和log三大类数据。

* metric数据一般主要含有时间、指标名称、指标值和指标属性等字段。特点是数据量与系统指标成正比，数据为典型的时序数据，适合列式存储，压缩性价比较高。
* trace数据为链式的数据，一个trace中含有多个span，span间存在关联关系，适用于系统调用状态的状态分析。同种trace类型下数据结构较为固定。数据生成时需要传递trace的父span信息用于调用链的后续调用。

* log数据与系统以及业务挂钩，具有数据结构复杂多变、数据量与系统业务量相关的特点。

在某些时间敏感的场景下，对日志输出模块的调用要求可能较高。数据生成模块在数据对数据生成具有一定的连续性，从而形成数据上下文，那么数据输出模块应保持这种数据上下文关系，避免数据乱序。

总的说数据输出模块应向上提供灵活的数据格式，保持数据上下文关系，性能上尽可能的不影响主业务运行，多数据输出场景下支持数据输出路由。

## 2.2 数据采集模块特性

* 日志数据采集模块实例与应用端日志数据输出实例，数量上往往为n:m。日志采集与日志输出模块应形成一定的划分规律
* 由于日志采集模块的单线程性能限制，日志采集模块可能会需要日志输出模块能够支持多个日志输出，使得日志采集模块可以并发的采集，从而提高日志采集的并发量，提高采集效率。
* 日志采集模块对日志的消费速度应大于日志输出产生速度，即日志数据不应产生堆积，否则堆积日志会因生命周期管理而造成丢失或占用大量额外存储。
* 且数据采集和数据输出就数据协议上应达成一致。并尽可能使用具有高压缩比，高性能的序列化协议，确保采集模块的初步解析性能。
* 日志采集模块往往还肩负日志数据存储库路由职责，因此日志采集模块往往需要对日志进行初步解析，获取相关元信息用于日志输出路由决策。在选择数据序列化协议时应注意该问题。

简单说应注意数据的输出位置、数据输出协议，数据元信息共识、数据类型共识、数据序列化协议。



# 三、数据输出模块需求分析

总结相邻模块的特性的相关需求：

1. 数据输出模块应向上提供灵活的数据结构表示log类型日志
2. 提供metric类型数据的输出
3. 提供Trace调用链生成和输出
4. 基本log日志行为支持
5. 一定程度的易用性
6. 保持数据上下文关系
7. 支持多输出数据路由
8. 数据元信息补充
9. 数据开放标准兼容
10. 数据的序列化
11. 序列化数据的输出
12. 不影响系统主业务
13. 其他



# 四、模块设计

由于业务的复杂与多变性，难以提供能够满足所有场景的数结构，因此采用分层的方式，在应用层的数据表示应有应用层转换为或使用日志输出模块提供的数据表现形式。数据输出模块提供若干个日志数据结构提供给应用层使用，这些数据结构实现了数据输出模块的数据表示层，该层会贯彻整个日志数据输出模块。如果有必要再对提供的数据结构类型进行扩展。



## 4.1 分层设计

### 4.1.2 分层策略

根据数据输出模块的需求，我们将各个需求进行分层，不同的需求在不同的层次满足其需求。整个分层结构如

1. 第一层数据输出层：负责将字节流输出到特定输出
2. 第二层数据序列化层：负责将上层传入的数据表示按要求进行序列化
3. 第三层数据开放协议层：负责将上层的传入的数据按数据开放标准进行转换与兼容，甚至添加一定的元信息
4. 第四层日志运行输出时：负责将上层传入的数据，根据数据元信息与一定的策略进行数据路由，选择数据输出，并保证数据上下文关系
5. 第五层数据表示层：实现对metric、trace、log类型数据的表示，并提供相应使用接口。
6. 第六层日志器操作层：提供基本的日志行为操作，负责对数据的日志级别过滤、采样率处理等，以及一定程度的易用性。
7. 第七层数据生成层：该层由上层应用系统实现，负责生成数据，并使用第五层提供的接口和数据结构表示，使用第六层进行数据记录。

分层结构如下图表示：

![](images/span%20layer%20design-span%20layer.png)

### 4.1.3 缺陷

数据表示层实际上不是严格意义上的独立层，它被第二到第七层所依赖。一旦未来对数据表示层进行重写，将会影响所有的分层。



### 4.1.4 优势

* 可以根据需要增加分层，提供更丰富的功能。
* 灵活的增加分层实现，有利于扩展分层能力
* 对某层的职责变更或接口改写，仅会影响相邻层，影响较小
* 上层仅需要关注邻层接口，不需要关心其他层的使用



## 4.2 对象设计

### 4.2.1 数据表示对象

![](images/span-span%20log.png)

#### 4.2.1.1 介绍

图中蓝色的interface表示概念对象，即trace、metric、log等类型日志。

* `Trace`：Trace代表链路追踪类型日志，如上类图所示Trace逻辑主要由`InternalSpan`接口实现，该接口代表了系统的内部调用。这么设计的原因是所有的外部调用追踪(例如http链路追踪)本质上还是由内部调用发起，虽然该内部调用并非整个调用链的头部。所以在实现上外部调用`externalSpan`与`InternalSpan`为聚合关系，应由`InternalSpan`提供创建接口。
* `InternalSpan`：代表的是系统内部调用的链路追踪，同时也是其他类型日志的容器。
  * 线程间不允许共有同一`InternalSpan`，具有调用链关系的`InternalSpan`应为父子或兄弟关系。注意该对象操作是**非线程安全**的，在当前线程使用结束后应通过`Signal()`接口释放
  * 通过该对象完成了系统的内部链路追踪和保证上下文连续性。其他类型的日志都应该由该对象进行分发或记录。
* `externalSpan`：`externalSpan`代表外部调用，典型代表为http链路追踪、kafka每个span。目前实现中通过`InternalSpan`创建对应对象，并由使用者通过对应类型提供的接口填充信息。应在其父`InternalSpan`释放前完成信息记录。
* `log`：通过一个`Field`实例完成对复杂日志事件的表示，并最终添加到当前线程的`InternalSpan`中，形成上下文，不同线程的`log`必须添加到各自的`InternalSpan`，以确保日志事件上下文的连续性，并且可以实现内部链路追踪信息丰富。
* `metrics`：通过`Mmetric`完成对指标类型日志的表示和记录。实现上需要添加到当前线程的`InternalSpan`中，虽然`metrics`不需要由`InternalSpan`保管，但是通过`InternalSpan`，可以实现批量记录，以及为内部调用提供丰富的数据。
* `Field`：主要用于对复杂多变的日志事件表示



#### 4.2.1.2 想法

* `InternalSpan`实现了系统内部线程链路追踪或者事务链路追踪
* `InternalSpan`不保证线程安全，但是每个线程都使用不同`InternalSpan`实例，因此不存在需要线程安全的情况。
* 子线程使用的`InternalSpan`由父线程创建后传递，同样不需要线程安全。
* log、metric必须使用`InternalSpan`作为容器，使得同一线程的内的日志上下文可以保持连续性，同时在需要系统内部链路追踪时可以得到丰富的信息。
* `InternalSpan`作为数据输出的最小单位，可以有效的减少零碎的数据记录行为，将其变为批处理，减少对数据输出的竞争，提高性能。



#### 4.2.1.3 缺点

* 使用时需要传递`InternalSpan`，比起一般日志器增加了额外条件
  * 考虑通过封装接口，结合context完成对额外条件的掩盖



### 4.2.2 运行时与分层对象

![](./images/span-span%20runtime.png)

#### 4.2.2.1 介绍

* `Runtime`：提供日志运行时，对应日志运行时分层，每个祖线程的`InternalSpan`应该由`Runtime`提供，在`InternalSpan`释放后会交由运行时处理。完成数据路由，并保证该层数据操作的线程安全。
* `Encoder`：对应数据开放协议兼容层，负责将`InternalSpan`代表数据按要求进行填充信息、转换、兼容，并调用序列化器完成序列化
* `FieldWriter`：对象数据序列化层，将`Field`对象序列化为字节流，并将字节流传入输出层。
* `Writer`：输出层，负责将数据写入磁盘、网络等。甚至可以写入缓存等，形成嵌套，目前在实现上就是一个`io.Writer`



#### 4.2.2.2 想法

* `Writer`完成字节流的输出，目前实现中会通过将缓存与文件输出组合降低缓存到文件的写性能差异。
* `FieldWriter`直接将字节流通过输出层接口写入输出层，能降低字节数组的合并行为，如果在实现上多次进行字节数组的合并，则容易发生数组重分配，造成不必要的数据复制。且`FieldWriter`序列化所需的缓存数组也可以由自身有效的控制。
* `Encoder`来完成数据开放协议标准兼容，这样使得序列化层不需要考虑标准问题，当变更标准或序列化协议时两者不会有太大的影响。一个`FieldWriter`和仅持有一个`writer`
* `Runtime`：负责系统内每个事务第一个线程`InternalSpan`的分发，并接收释放的`InternalSpan`，在合适时机对其进行数据路由，选择对应的`Encoder`进行编码与输出



#### 4.2.2.3 缺点

* `Writer`和`FieldWriter`绑定，这意味着该线程要么在序列化，要么在进行IO输出，这两个行为一个是IO行为一个是cpu计算，两种行为没有很好的进行分离。如果输出层写入性能较慢，则会造成该线程大部分时间处于io阻塞状态，不能充分利用cpu。目前通过内存缓存提高了一部分性能，但是仍然存在该问题。
  * 可以考虑进一步处理输出层，使得整个流程变为变序列化边输出，将序列化和输出放在不同的线程上执行。但是线程间同步会带来额外开销，以目前的实现在整个性能测试的表现中来看，IO等待时间为0.24s，而线程间runtime的数据同步确需要0.27s。因此该方案是否可行仍需要考虑，并进一步测试验证
* `Runtime`关闭时机，日志器应该是整个系统最后关闭的模块，通过系统信号量的方式触发`Runtime`关闭则有可能在其他模块需要记录日志时就已经关闭。因此需要使用这手工关闭。
* 整个日志器初始化构造过程可能会比较复杂。
* 有效需求下`FieldWriter`和`Writer`并不一定能够完全分离，例如当`Writer`不时字节流输出，而是对象输出时，两者可能需要融合



## 4.3 使用用例

### 4.3.1 用例图

![](images/span%20layer%20design-span%20usecase.png)



### 4.3.2 example

例子中包含单线程日志数据记录和多线程下数据记录的操作。需要注意的是`InternalSpan`的操作是非线程安全的，并且在使用结束后必须要通过`Signal`手工释放

```go
import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"span/encoder"
	"span/field"
	"span/open_standard"
	"span/runtime"
	"testing"
	"time"

	"gotest.tools/assert"
)


func TestSamplerLogger(t *testing.T) {
	// 0. create logger and start runtime
	buf := bytes.NewBuffer(nil)
	l := NewdefaultSamplerLogger()
	run := runtime.NewRuntime(&open_standard.OpenTelemetry{
		Encoder: encoder.NewJsonEncoder(buf),
	}, field.NewSpanFromPool)
	l.SetRuntime(run)
	l.LogLevel = AllLevel
	go run.Run()

	// 1. first create a root internalSpan
	root := l.NewInternalSpan()

	// 1.0 set trace info for root internalSpan
	traceID := field.GenSpanID()
	externalParentID := field.GenSpanID()
	l.SetTraceID(traceID, root)
	l.SetParentID(externalParentID, root)

	// 1.1 log message into roor internalSpan
	l.Debug("debug string message", root)
	l.DebugField(field.StringField("debug field message"), root)

	// 1.2 create a child internalSpan from root for a sub thread/task
	child0 := l.ChildrenInternalSpan(root)

	// 1.3 start a new thread for sub task
	go func() {
		// 2.1 log message into child internalSpan for child thread
		l.Debug("debug string", child0)

		// 2.X signal child0
		child0.Signal()
	}()

	// 1.4 record some metric into root internalSpan
	m := field.Mmetric{}
	m.Set("root thread", 0.0)
	m.AddLabel("root")
	m.AddLabel("metric")
	m.AddAttribute("root", "root span")
	l.RecordMetrics(m, root)

	// 1.5 record first external request into root internalSpan
	es, err := l.NewExternalSpan(root)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// 1.5.1 get trace info for some work
	tID := es.TraceID()
	espID := es.ParentID()
	parentID := es.ParentID()
	spanID := es.ID()
	// 1.5.2 write info to external span
	es.StartTime = time.Now()
	es.EndTime = time.Now()
	es.Attributes.Set("method", field.StringField("test"))
	es.Attributes.Set("host", field.StringField("test"))
	es.Attributes.Set("attr0", field.StringField(tID))
	es.Attributes.Set("attr1", field.StringField(espID))
	es.Attributes.Set("attr2", field.StringField(parentID))
	es.Attributes.Set("attr3", field.StringField(spanID))

	// 1.X signal root internalSpan
	root.Signal()

	// final close runtime and clean work space
	l.Close()
	// run.Signal()

	// check test result
	assert.Equal(t, traceID, tID)
	assert.Equal(t, externalParentID, espID)

	cap := map[string]interface{}{}
	bytes := buf.Bytes()
	left := 0
	i := 0
	n := 0
	for ; i < len(bytes); i += 1 {
		if bytes[i] == '\n' {
			if err = json.Unmarshal(bytes[left:i], &cap); err != nil {
				t.Error(err)
				t.FailNow()
			} else {
				n += 1
				fmt.Println(string(bytes[left:i]))
				fmt.Println()
			}
			left = i + 1
		}
	}
	if left < len(bytes) {
		if err = json.Unmarshal(bytes[left:i], &cap); err != nil {
			t.Error(err)
			t.FailNow()
		} else {
			n += 1
			fmt.Println(string(bytes[left:i]))
		}
	}

	// fmt.Print(buf.String())
}
```





# 五、实现

主要介绍各个分层实现选择与原因以及存在的问题

## 5.1 输出层

输出层采用内存缓存+标准输出的方式作为输出

### 5.1.1 优点

* 使用标准输出的好处在于简单，在容器平台上的标准输出会被重定向到特定文件，不需要考虑日志临时存储问题。
* 使用内存缓存可以降低外部存储与内存的io性能差，减少系统调用，增加IO性能
* 一个输出对应一个缓存，使用的缓存最大为`Field`中最大字符串大小。

### 5.1.2 缺点

* 不支持多输出，单输出成为瓶颈，在性能要求较高的情况下可能难以满足上层应用性能要求
* 不支持多输出，采集端无法在数据读取部分进行并发操作
* 与其他第三方模块不兼容，因为其他第三方模块不会使用自研日志输出模块

### 5.1.3 其他方案

1. 多文件输出

   * 优点在于支持了多输出能力，提高并发度
   * 相对复杂
     1. 需要运行时支持定义多输出路由策略。
     2. 文件输出着需要挂载存储作为日志临时存储方案。这使得容器平台上的部署复杂度增加，解决办法是通过提供统一的日志存储卷使用策略，以及部署模板注入方式降低部署复杂度，使得挂载存储过程透明化

2. 网络输出

   * 优点在于不需要本地存储
   * 缺点在于IO受限于日志接收器，可能性能问题。并且可能由于网络问题导致数据丢失等。

   

## 5.2 序列化层

目前使用json作为序列化协议，并且json序列化不是使用标准的go encoding/json包进行编码。

自行编写方法基于`Field`对象，完成json序列化。在写字符串时，需要转换为字节数组，此时基于非安全的转换形式转换为字节数组，底层共享同一片内存。非标准的字符串到字节的转换，避免了额外的内存申请和内存复制，性能更好。

### 5.2.1 优点

1. json序列化不依靠标准库，而是自行编写，序列化性能更好。encoding库依赖于go的反射机制，性能较差。基于`Field`对象可以避免这部分开销，而且可以将字节数据传递到输出，无需额外的字节数组内存创建。
2. json序列化比较通用，且便于人类阅读。

### 5.2.2 缺点

1. 序列化性能较差

   1. 传入的字符串需要进行额外的安全编码转换
   2. 本质上是字符串，类型描述上比较模糊。例如只能知道是数字而无法区分数字类型。因此传入的类型需要二次转换
   3. 冗余字符较多，例如冒号

2. 反序列化性能差，而且需要序列化整个结构，在一部分场景下仅需要部分序列化

3. 缺少基本的协议相关元信息，不利于采集端发现数据异常、快速路由和数据传递

4. 获取数据元信息需要解析整个json字符串，不利于采集端快速获取需要的元信息和类型信息，进行初步解析和快速数据路由

5. 数据体积较大

   

### 5.2.3 其他方案

1. 使用msgpack，优点是可以减少特殊的字符，支持局部解析，减少体积，缺点是不利于人类阅读
2. gopobuf，优点是数据表示更好的降低数据体积，也是opentelemetry标准使用。缺点是goprobuf不包含数据key值信息，需要输出和接收方就数据结构达成一致，否则无法解析，灵活性较低。不利于人类阅读。
3. 自定义序列化协议，更贴合采集需求，支持局部解析，有利于元数据信息提供，协议版本控制，与目前的`Filed`设计是相符合的。



## 5.3 标准转换层

将日志数据按数据开放标准进行转换，目前使用opentelemetry数据定义。

### 5.3.1 优点

1. 大家都在吹，所以也许会很通用

### 5.3.2 缺点

1. 似乎缺少协议元信息定义，例如版本



## 5.4 运行时层

负责完成对`InternalSpan`的初步分发，与最终回收。运行时运行在独立的线程下，当其分发的`InternalSpan`释放时将会通过go chan将自身传递给运行时，运行时将根据`InternalSpan`进行数据路由。目前没有数据路由能力。考虑后续添加。

### 5.4.1 考量

1. 让运行时在独立的线程，通过管道回收`InternalSpan`，实现日志的异步处理，第一目的在于尽可能的降低对上层应用在写日志时对应用线程的开销，降低对时间敏感业务的影响
2. 使应用可以异步处理日志数据，另一个目的将日志数据处理的线程安全问题交给运行时。上层应用实际上并不需要关注日志数据是否确实罗盘，仅需要保证日志已经被下层接收。

### 5.4.2 缺点

1. 上层应用线程将数据同步到运行时线程使用了go的chan管道，带来额外开销



## 5.5 数据表示层

目前日志事件的表示是通过`Field`表示，指标通过`Mmetric`结构，`Trace`通过`InternalSpan`分发，应用进行填写。上层应用在事务开始时应从运行时申请第一个`InternalSpan`，在创建其他子线程或需要记录子调用/子任务时前，应通过当前线程所有的`InternalSpan`创建子`InternalSpan`，并交给子线程使用。在当前线程结束时需要手工释放当前线程的`InternalSpan`。

每个`InternalSpan`释放会创建独立的线程进行，不阻塞上层应用。`InternalSpan`释放会等待子`InternalSpan`的释放，直到所有子`InternalSpan`释放后才会释放自身。根`InternalSpan`，即运行时创建的`InternalSpan`释放时会通过管道将自身所有权传递给运行时线程。

`InternalSpan`分发的`ExternalSpan`不是真正的span，而是一个外部调用表示的指针。因为外部调记录仅需要仅记录相关信息，不会再次在系统内部生成子span，其子span由外部系统完成记录。如果外部调用在需要在子线程记录，应创建子`InternalSpan`用于记录。

### 5.5.1 考量

1. `InternalSpan`对象不是线程安全的，因为在使用上每个线程都应该有自身独立的`InternalSpan`，所以不需要线程安全。`InternalSpan`带来的额外好处是能够实现系统内部线程/任务的链路调用追踪
2. `InternalSpan`释放等待子`InternalSpan`是希望能够将整个内部调用形成的Trace打包再一块，保持调用上下文的连续性，减少数据分析需要额外进行的聚合。同时在数据量较大时，可以降低对go chan管道的压力。如果确实需要快速的打印span，而不需要等待子span，可以进行调整，让一个Trace的各个span独立输出。
3. 日志事件通过`Field`表示提供一定的日志记录灵活性
4. `InternalSpan`作为日志事件的容器，将当前线程的所有日志一起记录，目的在于保证线程上下文的连续性。因为线程离散的日志上下文不利于分析，而且有利于丰富系统调用分析上下文。metric因为相同的原因使用`InternalSpan`作为容器

## 5.5.2 缺点

1. 序列化层至应用层都依赖于`Field`，并非只有邻层依赖。所以准确的说该层不能算合理的独立层，分离得不够好，一旦进行接口调整，可能会影响所有，可以再考虑考虑。
2. 在上层应用记录日志时，由于使用`InternalSpan`可能会造成日志没有输出的错觉，实际上是在等待整个内部调用结束。这个缺点可以考虑将Trace的各个`InternalSpan`独立输出解决。日志事件的输出需要在`InternalSpan`处理时输出，即需要等待该线程结束，在开放debug时会给上层应用没有输出的错觉，有时开放排错会觉得不够方便，这个可以通过使用日志器的封装的易用接口立即输出日志事件解决该问题。
3. `InternalSpan`需要在显式传递和使用，使用不够便捷，考虑与context结合，封装接口



## 5.6 日志器层

日志器层主要完成日志级别过滤，和日志采样率的功能。同时会为应用层提供日志操作接口。提供一定程度的易用性接口。

目前记录日志时，日志器需要传入`InternalSpan`和对应的日志数据



### 5.6.1 考量

1. 需要传入`InternalSpan`是因为日志器主要做一些日志相关业务处理，但是不保证日志上下文连续性，需要`InternalSpan`
2. 当传入`InternalSpan`为空时，日志器将会自动从运行时申请`InternalSpan`在记录对应数据后立即释放，这个可以提供一定程度的易用性，同时在开放进行debug时需要立即进行日志输出时可以这么使用，但是这样将会丢失线程日志连续性的好处
3. 当传入的`InternalSpan`非空时，仅会根据一定策略进行日志记录，不会释放`InternalSpan`，对应`InternalSpan`释放仍由使用者释放。



### 5.6.2 缺点

1. 日志器的使用，需要传入`InternalSpan`，在使用上不够简单，可能后续会考虑结合context来进行简化



## 5.7 应用层

应用层根据自身需要使用`Field`对象表示数据，并使用日志器进行日志记录。该层由上层应用自行进行开发、细化分层。例如将审计日志使用

