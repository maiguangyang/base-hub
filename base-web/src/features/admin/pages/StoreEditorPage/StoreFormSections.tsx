import { Input } from '@/components/ui/input';
import { mobilePhoneInputProps } from '@/lib/mobilePhoneInput';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Checkbox } from '@/components/ui/checkbox';
import { TimeRangePicker } from '@/components/ui/time-picker';
import type { StoreBusinessStatus } from '@/__generated__/graphql';
import type { StoreValues } from './storeInputs';

interface SectionProps {
  values: StoreValues;
  change<K extends keyof StoreValues>(field: K, value: StoreValues[K]): void;
}

export function StoreBasicFields({ values, change }: SectionProps) {
  return (
    <div className="flex flex-col gap-3">
      <div className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">基础与联系</div>
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
        <label className="flex flex-col gap-2 text-sm sm:col-span-2">
          <span>门店名称 <span className="text-destructive">*</span></span>
          <Input placeholder="请输入门店名称，如：江南西总店" value={values.name} onChange={(e) => change('name', e.target.value)} maxLength={128} required />
        </label>
        <label className="flex flex-col gap-2 text-sm">
          <span>联系电话 <span className="text-destructive">*</span></span>
          <Input placeholder="请输入门店联系电话" value={values.contactPhone ?? ''} onChange={(e) => change('contactPhone', e.target.value)} maxLength={32} required />
        </label>
        <label className="flex flex-col gap-2 text-sm">
          <span>营业时间 <span className="text-destructive">*</span></span>
          <TimeRangePicker
            value={values.businessHours ?? '09:00 - 22:00'}
            onChange={(val) => change('businessHours', val)}
            placeholder="请选择营业时间（如 09:00 - 22:00）"
          />
        </label>
        <label className="flex flex-col gap-2 text-sm">
          <span>店长姓名</span>
          <Input placeholder="请输入店长姓名（选填）" value={values.managerName ?? ''} onChange={(e) => change('managerName', e.target.value)} maxLength={64} />
        </label>
        <label className="flex flex-col gap-2 text-sm">
          <span>店长手机号</span>
          <Input placeholder="请输入店长手机号（选填）" value={values.managerPhone ?? ''} onChange={(e) => change('managerPhone', e.target.value)} {...mobilePhoneInputProps} />
        </label>
      </div>
    </div>
  );
}

export function StoreLocationFields({ values, change }: SectionProps) {
  return (
    <div className="flex flex-col gap-3">
      <div className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">地理位置</div>
      <div className="grid grid-cols-1 gap-4 md:grid-cols-3">
        <label className="flex flex-col gap-2 text-sm">
          <span>省份 <span className="text-destructive">*</span></span>
          <Input placeholder="如：广东省" value={values.province ?? ''} onChange={(e) => change('province', e.target.value)} maxLength={64} required />
        </label>
        <label className="flex flex-col gap-2 text-sm">
          <span>城市 <span className="text-destructive">*</span></span>
          <Input placeholder="如：广州市" value={values.city ?? ''} onChange={(e) => change('city', e.target.value)} maxLength={64} required />
        </label>
        <label className="flex flex-col gap-2 text-sm">
          <span>区县 <span className="text-destructive">*</span></span>
          <Input placeholder="如：天河区" value={values.district ?? ''} onChange={(e) => change('district', e.target.value)} maxLength={64} required />
        </label>
      </div>
      <label className="flex flex-col gap-2 text-sm">
        <span>详细地址 <span className="text-destructive">*</span></span>
        <Input placeholder="请输入详细门牌地址，如：天河路228号正佳广场4楼" value={values.address ?? ''} onChange={(e) => change('address', e.target.value)} maxLength={256} required />
      </label>
    </div>
  );
}

export function StoreOperationFields({ values, change }: SectionProps) {
  return (
    <div className="flex flex-col gap-3">
      <div className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">经营与服务模式</div>
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
        <label className="flex flex-col gap-2 text-sm">
          <span>营业状态</span>
          <Select
            value={values.businessStatus ?? 'OPEN'}
            onValueChange={(val) => change('businessStatus', val as StoreBusinessStatus)}
          >
            <SelectTrigger className="w-full">
              <SelectValue placeholder="请选择营业状态" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="OPEN">正常营业</SelectItem>
              <SelectItem value="CLOSED">暂停营业</SelectItem>
            </SelectContent>
          </Select>
        </label>
        <div className="flex flex-col gap-2 text-sm justify-center">
          <span>服务模式</span>
          <div className="flex items-center gap-6 h-9">
            <label className="flex items-center gap-2 cursor-pointer select-none">
              <Checkbox checked={values.supportDineIn ?? true} onCheckedChange={(checked) => change('supportDineIn', Boolean(checked))} />
              <span>支持堂食</span>
            </label>
            <label className="flex items-center gap-2 cursor-pointer select-none">
              <Checkbox checked={values.supportTakeout ?? true} onCheckedChange={(checked) => change('supportTakeout', Boolean(checked))} />
              <span>支持外带</span>
            </label>
          </div>
        </div>
      </div>
    </div>
  );
}

export function StoreFacilityFields({ values, change }: SectionProps) {
  return (
    <div className="flex flex-col gap-3">
      <div className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">设施与小票（选填）</div>
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
        <label className="flex flex-col gap-2 text-sm">
          <span>经营面积（㎡）</span>
          <Input type="number" placeholder="请输入经营面积（选填）" value={values.storeArea ?? ''} onChange={(e) => change('storeArea', e.target.value)} />
        </label>
        <label className="flex flex-col gap-2 text-sm">
          <span>桌位数（桌）</span>
          <Input type="number" placeholder="请输入桌位数（选填）" value={values.tableCount ?? ''} onChange={(e) => change('tableCount', e.target.value)} />
        </label>
      </div>
      <label className="flex flex-col gap-2 text-sm">
        <span>小票底部寄语</span>
        <Input placeholder="例如：凭小票可向前台领取2小时免费停车券，欢迎再次光临！" value={values.receiptFooter ?? ''} onChange={(e) => change('receiptFooter', e.target.value)} maxLength={256} />
      </label>
    </div>
  );
}
