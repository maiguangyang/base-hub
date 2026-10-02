# 自动化切图与光栅化渲染脚本规范 (Swift & AppKit)

在 macOS 环境下，利用原生 Swift + AppKit / CoreGraphics 引擎可以实现毫秒级、免外部庞大依赖的高保真切图与矢量渲染。

---

## 1. 沙箱执行命令规范

在 Antigravity / Codex 沙箱内执行 Swift 脚本时，必须显式指定模块缓存路径（避免因写入系统目录报错）：

```bash
swift -module-cache-path /Users/marlon.m/.gemini/antigravity/brain/<conv-id>/scratch/cache <script_path>.swift
```

---

## 2. 核心 AppKit 位图创建与防双倍缩放模版

```swift
import AppKit
import CoreGraphics

/// 创建标准 8-bit RGBA 位图上下文
func createBitmap(width: Int, height: Int) -> (NSBitmapImageRep, CGContext) {
    let rep = NSBitmapImageRep(
        bitmapDataPlanes: nil,
        pixelsWide: width,
        pixelsHigh: height,
        bitsPerSample: 8,
        samplesPerPixel: 4,
        hasAlpha: true,
        isPlanar: false,
        colorSpaceName: .deviceRGB,
        bytesPerRow: 0,
        bitsPerPixel: 0
    )!
    // 关键防坑：rep.size 必须与 pixelsWide/pixelsHigh 一致，避免绘制时产生 2x 错位或变形
    rep.size = NSSize(width: width, height: height)
    NSGraphicsContext.current = NSGraphicsContext(bitmapImageRep: rep)
    let ctx = NSGraphicsContext.current!.cgContext
    return (rep, ctx)
}

/// 保存为高质量 PNG
func savePNG(rep: NSBitmapImageRep, to path: String) {
    let data = rep.representation(using: .png, properties: [:])!
    try! data.write(to: URL(fileURLWithPath: path))
    print("Saved: \(path)")
}
```

---

## 3. 坐标系转换工具 (Top-Left ➔ AppKit Bottom-Left)

AppKit 的默认坐标系原点 `(0, 0)` 在左下角，而现代设计稿坐标均以左上角为原点。通过以下工具函数进行统一换算：

```swift
func toAppKitY(canvasH: Int, topY: CGFloat, itemH: CGFloat) -> CGFloat {
    return CGFloat(canvasH) - topY - itemH
}
```

---

## 4. 矢量 SVG 直接高质量渲染为 192px PNG 模版

macOS 11+ 原生 `NSImage(contentsOfFile:)` 支持无损加载并光栅化 SVG：

```swift
import AppKit

func renderSvgToPng(svgPath: String, pngPath: String, size: Int = 192) {
    guard let img = NSImage(contentsOfFile: svgPath) else {
        fatalError("Failed to load SVG: \(svgPath)")
    }
    let (rep, _) = createBitmap(width: size, height: size)
    img.draw(in: NSRect(x: 0, y: 0, width: size, height: size))
    savePNG(rep: rep, to: pngPath)
}
```

---

## 5. 合成全屏终端暗幕 Mockup 模版

```swift
import AppKit
import CoreGraphics

func compositeKioskMockup(
    bgMockupPath: String,
    dialogCardPath: String,
    outPath: String,
    screenW: Int = 768,
    screenH: Int = 1376,
    cardW: Int = 560,
    cardH: Int = 820
) {
    let (mockRep, mockCtx) = createBitmap(width: screenW, height: screenH)

    // 1. 绘制底层主屏 (如购物车或欢迎页)
    if let bgImg = NSImage(contentsOfFile: bgMockupPath) {
        bgImg.draw(in: NSRect(x: 0, y: 0, width: screenW, height: screenH))
    }

    // 2. 覆盖全屏暗幕遮罩 (rgba(15, 23, 42, 0.65))
    mockCtx.setFillColor(NSColor(red: 0.06, green: 0.09, blue: 0.16, alpha: 0.65).cgColor)
    mockCtx.fill(CGRect(x: 0, y: 0, width: screenW, height: screenH))

    // 3. 居中绘制弹窗卡片
    let cardX = (screenW - cardW) / 2
    let cardY = (screenH - cardH) / 2
    if let cardImg = NSImage(contentsOfFile: dialogCardPath) {
        cardImg.draw(in: NSRect(x: cardX, y: cardY, width: cardW, height: cardH))
    }

    savePNG(rep: mockRep, to: outPath)
}
```
