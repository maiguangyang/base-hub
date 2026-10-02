# Kiosk UI 资产库目录与分层架构标准

## 1. 根目录结构

`kiosk_ui_assets` 位于项目根目录，采用**「公共沉淀（Common）+ 按业务页面独立分包（Pages）」**的架构体系。

```text
kiosk_ui_assets/
├── index.html                               # ★ 全局可视化看板与全流程多页面 Tab 预览中心
├── README.md                                # 架构说明与资产树文档
│
├── common/                                  # 【多页面共享的全局公共资产】
│   ├── logos/                               # 品牌 Logo (字标、徽章、Slogan)
│   │   ├── logo_option1_calligraphy.svg/.png# 定稿国潮书法字标
│   │   └── logo_review.html                 # 品牌 Logo 比选看板
│   └── icons/                               # 全局系统级通用图标 (192×192 PNG + SVG)
│       ├── icon_language.svg/.png           # 多语言切换
│       ├── icon_help_staff.svg/.png         # 呼叫店员协助
│       ├── icon_barcode_scan.svg/.png       # 扫码导引
│       └── icon_contactless_pay.svg/.png    # NFC感应支付
│
└── pages/                                   # 【按业务页面独立分包】
    ├── 01_welcome_screen/                   # 01. 欢迎主屏幕
    ├── 02_scan_cart/                        # 02. 商品扫码加购页（含手动输码弹窗）
    ├── 03_member_login/                     # 03. 会员手机号极速登录页 / 弹窗
    ├── 04_payment_methods/                  # 04. 聚合支付多通道选择页
    └── 05_receipt_success/                  # 05. 支付成功与小票打印页
```

---

## 2. 单个页面的 4 个标准子目录规范

每一个 `pages/XX_page_name/` 目录下，必须包含且仅包含以下 4 个标准子目录：

| 子目录 | 作用与定位 | 典型尺寸与格式 | 命名规范 |
| :--- | :--- | :--- | :--- |
| `preview/` | **整体界面效果图与真机 Mockup**。<br>供看板展示、设计评审与交互走查。 | 全屏：`768×1376` PNG<br>弹窗卡片：`560×820` PNG (透明底) | `kiosk_xxx_mockup.png`<br>`dialog_xxx_card.png` |
| `components/` | **局部独立可抽离组件切图**。<br>卡片、输入框、键盘、商品切片等，背景透明。 | 宽度统一规范：`480px`<br>高度依内容自适应 (`76px` ~ `340px`) | `xxx_banner_card.png`<br>`xxx_formatted_display.png`<br>`xxx_tactile_pad.png` |
| `icons/` | **本页面专用的微图标**。<br>提供矢量源码与高质量光栅图。 | SVG 源码 (`viewBox="0 0 192 192"`)<br>高清 PNG (`192×192` 像素) | `icon_xxx.svg`<br>`icon_xxx.png` |
| `flutter/` | **生产级 Flutter Widget 源码**。<br>支持直接在 App 中一行代码呼出或嵌入页面。 | Dart 源码 (`.dart`) | `kiosk_xxx_screen.dart`<br>`dialog_xxx.dart` |

---

## 3. 镜像同步至 `base-app` 的路径映射

为确保客户端工程即开即用，新生成的资产必须同步镜像至 Flutter 项目目录：

- `pages/01_welcome_screen/` ➔ `base-app/assets/images/kiosk/pages/01_welcome/`
- `pages/02_scan_cart/` ➔ `base-app/assets/images/kiosk/pages/02_scan/`
- `pages/03_member_login/` ➔ `base-app/assets/images/kiosk/pages/03_member/`
- `pages/04_payment_methods/` ➔ `base-app/assets/images/kiosk/pages/04_payment/`
- `pages/05_receipt_success/` ➔ `base-app/assets/images/kiosk/pages/05_receipt/`

每次新增页面后，必须在 `base-app/pubspec.yaml` 中登记新增的 `preview/`、`components/`、`icons/` 路径。
