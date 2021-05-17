// Copyright (c) 2016 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package benchmarks

import (
    "bufio"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os"
    "span/encoder"
    "span/field"
    "span/open_standard"
    "span/runtime"
    "sync"
    "testing"

    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

func BenchmarkEncodeAndMalloc(b *testing.B) {
    b.Log("BenchmarkEncodeAndMalloc")

    b.Run("json/object", func(b *testing.B) {
        tmp := ioutil.Discard
        b.ResetTimer()
        // r := fakeSpanStructField()
        b.RunParallel(func(pb *testing.PB) {
            for pb.Next() {
                r := fakeMap()
                bytes, _ := json.Marshal(r)
                tmp.Write(bytes)
            }
        })
    })

    b.Run("json/fakeStruct", func(b *testing.B) {
        tmp := ioutil.Discard
        b.ResetTimer()
        // r := fakeSpanStructField()
        b.RunParallel(func(pb *testing.PB) {
            for pb.Next() {
                r := fakeStruct()
                bytes, _ := json.Marshal(r)
                tmp.Write(bytes)
            }
        })
    })

    b.Run("Myspan/object", func(b *testing.B) {
        // writer := &open_standard.OpenTelemetry{
        //  Encoder: encoder.NewJsonEncoderBench(ioutil.Discard),
        // }
        // logger := runtime.NewRuntime(writer, field.NewSpanFromPool)
        // go logger.Run()
        // defer logger.Signal()

        enc := encoder.NewJsonEncoderBench(ioutil.Discard)
        b.ResetTimer()
        b.RunParallel(func(pb *testing.PB) {
            for pb.Next() {
                r := fakeSpanStructField()
                enc.Write(r)
            }
        })

    })

    b.Run("Zap.Check", func(b *testing.B) {
        logger := newZapLogger(zap.DebugLevel)
        b.ResetTimer()
        b.RunParallel(func(pb *testing.PB) {
            for pb.Next() {
                if ce := logger.Check(zap.InfoLevel, getMessage(0)); ce != nil {
                    ce.Write(fakeFields()...)
                }
            }
        })
    })

}

func BenchmarkEncodeDiscardWrite(b *testing.B) {
    b.Log("BenchmarkEncodeDiscardWrite")

    b.Run("json/object", func(b *testing.B) {
        r := fakeMap()
        tmp := ioutil.Discard
        b.ResetTimer()
        // r := fakeSpanStructField()
        b.RunParallel(func(pb *testing.PB) {
            for pb.Next() {
                bytes, _ := json.Marshal(r)
                tmp.Write(bytes)
            }
        })
    })

    b.Run("json/fakeStruct", func(b *testing.B) {
        r := fakeStruct()
        tmp := ioutil.Discard
        b.ResetTimer()
        // r := fakeSpanStructField()
        b.RunParallel(func(pb *testing.PB) {
            for pb.Next() {
                bytes, _ := json.Marshal(r)
                tmp.Write(bytes)
            }
        })
    })

    b.Run("Myspan/object", func(b *testing.B) {
        r := fakeSpanStructField()
        enc := encoder.NewJsonEncoderBench(ioutil.Discard)
        b.ResetTimer()
        b.RunParallel(func(pb *testing.PB) {
            for pb.Next() {
                enc.Write(r)
            }
        })
    })

    b.Run("Zap.Check", func(b *testing.B) {
        // logger := newZapLogger(zap.DebugLevel)
        fields := fakeFields()
        tmp := zapcore.AddSync(ioutil.Discard)
        ec := zap.NewProductionEncoderConfig()
        ec.EncodeDuration = zapcore.NanosDurationEncoder
        ec.EncodeTime = zapcore.EpochNanosTimeEncoder
        enc := zapcore.NewJSONEncoder(ec)
        logger := zap.New(zapcore.NewCore(
            enc,
            tmp,
            zap.DebugLevel,
        ))
        b.ResetTimer()
        b.RunParallel(func(pb *testing.PB) {
            for pb.Next() {
                if ce := logger.Check(zap.InfoLevel, getMessage(0)); ce != nil {
                    ce.Write(fields...)
                }
            }
        })
    })
}

func BenchmarkDiscardWrite(b *testing.B) {
    b.Log("BenchmarkDiscardWrite")

    b.Run("json/object", func(b *testing.B) {
        tmp := ioutil.Discard
        locker := &sync.Mutex{}
        b.ResetTimer()
        // r := fakeSpanStructField()
        b.RunParallel(func(pb *testing.PB) {
            for pb.Next() {
                r := fakeMap()
                bytes, _ := json.Marshal(r)
                locker.Lock()
                tmp.Write(bytes)
                locker.Unlock()
            }
        })
    })

    b.Run("Myspan/runtime", func(b *testing.B) {
        logger := runtime.NewRuntime(&open_standard.OpenTelemetry{
            Encoder: encoder.NewJsonEncoder(ioutil.Discard),
        }, field.NewSpanFromPool)
        go logger.Run()
        defer logger.Signal()
        b.ResetTimer()
        b.RunParallel(func(pb *testing.PB) {
            for pb.Next() {
                r := fakeSpanStructField()
                l := logger.Children()
                l.Record(r)
                l.Signal()
            }
        })
    })

    b.Run("Myspan/object", func(b *testing.B) {
        logger := runtime.NewRuntime(&open_standard.OpenTelemetry{
            Encoder: encoder.NewJsonEncoder(ioutil.Discard),
        }, field.NewSpanFromPool)
        go logger.Run()
        defer logger.Signal()

        b.ResetTimer()
        b.RunParallel(func(pb *testing.PB) {
            for pb.Next() {
                r := fakeSpanStructField()
                l := logger.Children()
                l.Record(r)
                l.Signal()
            }
        })
    })

    b.Run("Zap.Check", func(b *testing.B) {
        tmp := zapcore.AddSync(ioutil.Discard)

        ec := zap.NewProductionEncoderConfig()
        ec.EncodeDuration = zapcore.NanosDurationEncoder
        ec.EncodeTime = zapcore.EpochNanosTimeEncoder
        enc := zapcore.NewJSONEncoder(ec)
        logger := zap.New(zapcore.NewCore(
            enc,
            tmp,
            zap.DebugLevel,
        ))
        b.ResetTimer()
        b.RunParallel(func(pb *testing.PB) {
            for pb.Next() {
                fields := fakeFields()
                if ce := logger.Check(zap.InfoLevel, getMessage(0)); ce != nil {
                    ce.Write(fields...)
                }
            }
        })
    })
}

func BenchmarkEncodeFileWrite(b *testing.B) {
    b.Log("BenchmarkEncodeFileWrite")

    b.Run("json/object", func(b *testing.B) {
        r := fakeMap()
        tmp, err := ioutil.TempFile("", "zaptest")
        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }
        tmp1 := bufio.NewWriter(tmp)
        locker := &sync.Mutex{}
        b.ResetTimer()
        // r := fakeSpanStructField()
        b.RunParallel(func(pb *testing.PB) {
            for pb.Next() {
                bytes, _ := json.Marshal(r)
                locker.Lock()
                tmp1.Write(bytes)
                locker.Unlock()
            }
        })
    })

    b.Run("Myspan/runtime", func(b *testing.B) {
        // tmp := &writer.JsonWriter{
        //  W: ioutil.Discard,
        // }

        r := fakeSpanStructField()
        logger := runtime.NewRuntime(&open_standard.OpenTelemetry{
            Encoder: encoder.NewJsonEncoder(ioutil.Discard),
        }, field.NewSpanFromPool)
        go logger.Run()
        defer logger.Signal()
        b.ResetTimer()
        b.RunParallel(func(pb *testing.PB) {
            for pb.Next() {
                // r.Write(tmp)
                l := logger.Children()
                l.Record(r)
                l.Signal()
            }
        })
    })

    b.Run("Myspan/object", func(b *testing.B) {
        tmp, err := ioutil.TempFile("", "logger")
        if err != nil {
            fmt.Println(err)
        }
        tmp1 := bufio.NewWriter(tmp)
        logger := runtime.NewRuntime(&open_standard.OpenTelemetry{
            Encoder: encoder.NewJsonEncoder(tmp1),
        }, field.NewSpanFromPool)
        go logger.Run()
        defer logger.Signal()
        defer tmp.Close()

        r := fakeSpanStructField()
        b.ResetTimer()
        b.RunParallel(func(pb *testing.PB) {
            for pb.Next() {
                // r.Write(w)
                l := logger.Children()
                l.Record(r)
                l.Signal()
            }
        })
    })

    b.Run("Zap.Check", func(b *testing.B) {
        // logger := newZapLogger(zap.DebugLevel)
        fields := fakeFields()
        tmp, err := ioutil.TempFile("", "zaptest")
        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }
        ec := zap.NewProductionEncoderConfig()
        ec.EncodeDuration = zapcore.NanosDurationEncoder
        ec.EncodeTime = zapcore.EpochNanosTimeEncoder
        enc := zapcore.NewJSONEncoder(ec)
        logger := zap.New(zapcore.NewCore(
            enc,
            tmp,
            zap.DebugLevel,
        ))
        b.ResetTimer()
        b.RunParallel(func(pb *testing.PB) {
            for pb.Next() {
                if ce := logger.Check(zap.InfoLevel, getMessage(0)); ce != nil {
                    ce.Write(fields...)
                }
            }
        })
        tmp.Close()
    })
}

func BenchmarkFileWrite(b *testing.B) {
    b.Log("BenchmarkFileWrite")

    b.Run("json/object", func(b *testing.B) {
        tmp, err := ioutil.TempFile("", "zaptest")
        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }
        tmp1 := bufio.NewWriter(tmp)
        locker := &sync.Mutex{}
        b.ResetTimer()
        // r := fakeSpanStructField()
        b.RunParallel(func(pb *testing.PB) {
            for pb.Next() {
                r := fakeMap()
                bytes, _ := json.Marshal(r)
                locker.Lock()
                tmp1.Write(bytes)
                locker.Unlock()
            }
        })
    })

    b.Run("Myspan/object", func(b *testing.B) {
        tmp, err := ioutil.TempFile("", "logger")
        if err != nil {
            fmt.Println(err)
        }
        // tmp1 := bufio.NewWriter(tmp)
        tmp1 := tmp
        logger := runtime.NewRuntime(&open_standard.OpenTelemetry{
            Encoder: encoder.NewJsonEncoder(tmp1),
        }, field.NewSpanFromPool)
        go logger.Run()
        defer logger.Signal()
        defer tmp.Close()

        b.ResetTimer()
        b.RunParallel(func(pb *testing.PB) {
            for pb.Next() {
                r := fakeSpanStructField()
                l := logger.Children()
                l.Record(r)
                l.Signal()
            }
        })
    })

    b.Run("Zap.Check", func(b *testing.B) {
        // logger := newZapLogger(zap.DebugLevel)
        tmp, err := ioutil.TempFile("", "zaptest")
        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }
        ec := zap.NewProductionEncoderConfig()
        ec.EncodeDuration = zapcore.NanosDurationEncoder
        ec.EncodeTime = zapcore.EpochNanosTimeEncoder
        enc := zapcore.NewJSONEncoder(ec)
        logger := zap.New(zapcore.NewCore(
            enc,
            tmp,
            zap.DebugLevel,
        ))
        b.ResetTimer()
        b.RunParallel(func(pb *testing.PB) {
            for pb.Next() {
                fields := fakeFields()
                if ce := logger.Check(zap.InfoLevel, getMessage(0)); ce != nil {
                    ce.Write(fields...)
                }
            }
        })
        tmp.Close()
    })
}

func BenchmarkMmetrics(b *testing.B) {
    b.Logf("BenchmarkMmetrics")

    b.Run("json/metric", func(b *testing.B) {
        mm := map[string]interface{}{
            "test": 0,
        }

        attrs := map[string]string{}

        for i := 0; i < 10; i++ {
            k := fmt.Sprint("attrs", i)
            v := fmt.Sprint("test", i)
            attrs[k] = v
        }
        // mm := &field.Mmetric{}
        // mm.Set("test", 0)
        // key := "test"
        // for i := 0; i < 100; i++ {
        //  mm.AddAttribute(key, key)
        // }
        w := ioutil.Discard
        b.ResetTimer()
        b.RunParallel(func(pb *testing.PB) {
            for pb.Next() {
                bytes, _ := json.Marshal(mm)
                w.Write(bytes)
            }
        })
    })

    b.Run("Zap/add", func(b *testing.B) {
        fakeField := []zap.Field{zap.Int("test", 0)}
        for i := 0; i < 100; i++ {
            fakeField = append(fakeField, zap.String("test", "test"))
        }
        logger := newZapLogger(zap.DebugLevel)
        b.ResetTimer()
        b.RunParallel(func(pb *testing.PB) {
            for pb.Next() {
                logger.Info(getMessage(0), fakeField...)
            }
        })
    })

    b.Run("metrics/self", func(b *testing.B) {
        mm := &field.Mmetric{}
        mm.Set("test", 0)
        key := "test"
        for i := 0; i < 100; i++ {
            mm.AddAttribute(key, key)
        }

        buf := ioutil.Discard
        enc := encoder.NewJsonEncoderBench(buf)
        b.ResetTimer()
        b.RunParallel(func(p *testing.PB) {
            for p.Next() {
                enc.Write(mm)
            }
        })
    })
}

