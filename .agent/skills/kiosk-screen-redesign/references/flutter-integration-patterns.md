# Kiosk Flutter 生产级组件架构规范

## 1. 弹窗组件通用架构标准

Kiosk 弹窗 Widget 统一设计为**自闭环、低耦合、开箱即用**的 StatefulWidget，并提供标准的静态 `show` 方法：

```dart
class DialogSample extends StatefulWidget {
  final void Function(String value)? onConfirm;
  final VoidCallback? onDismiss;
  final VoidCallback? onCallStaff;

  const DialogSample({
    super.key,
    this.onConfirm,
    this.onDismiss,
    this.onCallStaff,
  });

  /// 静态快速呼出入口
  static Future<String?> show(BuildContext context, {VoidCallback? onCallStaff}) {
    return showDialog<String>(
      context: context,
      barrierDismissible: true,
      barrierColor: const Color(0xFF0F172A).withOpacity(0.65), // 标准暗幕遮罩
      builder: (ctx) => DialogSample(
        onConfirm: (val) => Navigator.pop(ctx, val),
        onDismiss: () => Navigator.pop(ctx),
        onCallStaff: () {
          Navigator.pop(ctx);
          onCallStaff?.call();
        },
      ),
    );
  }

  @override
  State<DialogSample> createState() => _DialogSampleState();
}
```

---

## 2. 动态数字分段展示算法模版

### 手机号 3-4-4 分段：
```dart
String get formattedPhone {
  if (_digits.isEmpty) return '';
  if (_digits.length <= 3) return _digits;
  if (_digits.length <= 7) {
    return '${_digits.substring(0, 3)} ${_digits.substring(3)}';
  }
  return '${_digits.substring(0, 3)} ${_digits.substring(3, 7)} ${_digits.substring(7)}';
}
```

### 商品条码 3-4-6 分段：
```dart
String get formattedBarcode {
  if (_digits.isEmpty) return '';
  if (_digits.length <= 3) return _digits;
  if (_digits.length <= 7) {
    return '${_digits.substring(0, 3)} ${_digits.substring(3)}';
  }
  return '${_digits.substring(0, 3)} ${_digits.substring(3, 7)} ${_digits.substring(7)}';
}
```

---

## 3. 实体触觉按键构建模版

```dart
Widget buildKeyItem({
  required String label,
  required VoidCallback onTap,
  Color bgColor = Colors.white,
  Color borderColor = const Color(0xFFE2E8F0),
  Color textColor = const Color(0xFF0F172A),
  Widget? customChild,
}) {
  return InkWell(
    onTap: onTap,
    borderRadius: BorderRadius.circular(14),
    splashColor: const Color(0xFF2F54EB).withOpacity(0.12),
    child: Container(
      height: 64,
      decoration: BoxDecoration(
        color: bgColor,
        borderRadius: BorderRadius.circular(14),
        border: Border.all(color: borderColor, width: 1.2),
        boxShadow: const [
          BoxShadow(
            color: Color(0x0A000000),
            blurRadius: 6,
            offset: Offset(0, 3),
          ),
        ],
      ),
      alignment: Alignment.center,
      child: customChild ?? Text(
        label,
        style: TextStyle(fontSize: 26, fontWeight: FontWeight.w800, color: textColor),
      ),
    ),
  );
}
```

---

## 4. 状态驱动 CTA 按钮模版

```dart
Widget buildSubmitButton({
  required bool isValid,
  required VoidCallback onSubmit,
  required String activeText,
  required String disabledText,
}) {
  return Container(
    height: 58,
    decoration: BoxDecoration(
      gradient: isValid
          ? const LinearGradient(
              colors: [Color(0xFF2F54EB), Color(0xFF1D39C4)],
              begin: Alignment.topCenter,
              end: Alignment.bottomCenter,
            )
          : null,
      color: isValid ? null : const Color(0xFFE2E8F0),
      borderRadius: BorderRadius.circular(29),
      boxShadow: isValid
          ? [
              BoxShadow(
                color: const Color(0xFF2F54EB).withOpacity(0.38),
                blurRadius: 18,
                offset: const Offset(0, 8),
              ),
            ]
          : null,
    ),
    child: Material(
      color: Colors.transparent,
      child: InkWell(
        onTap: isValid ? onSubmit : null,
        borderRadius: BorderRadius.circular(29),
        child: Center(
          child: Text(
            isValid ? activeText : disabledText,
            style: TextStyle(
              fontSize: 17,
              fontWeight: FontWeight.w800,
              color: isValid ? Colors.white : const Color(0xFF94A3B8),
            ),
          ),
        ),
      ),
    ),
  );
}
```
