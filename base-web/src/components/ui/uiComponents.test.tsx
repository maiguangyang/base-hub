import { renderToStaticMarkup } from 'react-dom/server';
import { describe, expect, it } from 'vitest';
import { Button } from './button';
import { Card, CardHeader, CardTitle, CardContent, CardFooter } from './card';
import { Label } from './label';
import { Switch } from './switch';
import { Tabs, TabsList, TabsTrigger, TabsContent } from './tabs';
import { RadioGroup, RadioGroupItem } from './radio-group';
import { SegmentedControl, SegmentedControlItem } from './segmented-control';
import { Pagination, PaginationContent, PaginationItem, PaginationLink } from './pagination';
import { OverflowTooltipText } from './overflow-tooltip-text';

describe('UI Primitives Rendering', () => {
  it('renders Button with variants correctly', () => {
    const defaultBtn = renderToStaticMarkup(<Button>确认</Button>);
    expect(defaultBtn).toContain('bg-primary');

    const destructiveBtn = renderToStaticMarkup(<Button variant="destructive">删除</Button>);
    expect(destructiveBtn).toContain('bg-destructive');
  });

  it('renders Card hierarchy correctly', () => {
    const cardMarkup = renderToStaticMarkup(
      <Card>
        <CardHeader>
          <CardTitle>卡片标题</CardTitle>
        </CardHeader>
        <CardContent>卡片正文</CardContent>
        <CardFooter>页脚</CardFooter>
      </Card>,
    );
    expect(cardMarkup).toContain('卡片标题');
    expect(cardMarkup).toContain('卡片正文');
    expect(cardMarkup).toContain('页脚');
    expect(cardMarkup).toContain('data-slot="card"');
  });

  it('renders Label correctly', () => {
    const markup = renderToStaticMarkup(<Label htmlFor="test">测试标签</Label>);
    expect(markup).toContain('测试标签');
    expect(markup).toContain('for="test"');
  });

  it('renders Switch correctly', () => {
    const markup = renderToStaticMarkup(<Switch aria-label="开关" />);
    expect(markup).toContain('data-slot="switch"');
    expect(markup).toContain('data-slot="switch-thumb"');
  });

  it('renders Tabs correctly', () => {
    const markup = renderToStaticMarkup(
      <Tabs defaultValue="tab1">
        <TabsList>
          <TabsTrigger value="tab1">选项卡1</TabsTrigger>
        </TabsList>
        <TabsContent value="tab1">内容1</TabsContent>
      </Tabs>,
    );
    expect(markup).toContain('选项卡1');
    expect(markup).toContain('data-slot="tabs"');
  });

  it('renders RadioGroup correctly', () => {
    const markup = renderToStaticMarkup(
      <RadioGroup defaultValue="1">
        <RadioGroupItem value="1" aria-label="单选1" />
      </RadioGroup>,
    );
    expect(markup).toContain('data-slot="radio-group"');
    expect(markup).toContain('data-slot="radio-group-item"');
  });

  it('renders SegmentedControl correctly', () => {
    const markup = renderToStaticMarkup(
      <SegmentedControl defaultValue="day">
        <SegmentedControlItem value="day">按日</SegmentedControlItem>
        <SegmentedControlItem value="month">按月</SegmentedControlItem>
      </SegmentedControl>,
    );
    expect(markup).toContain('按日');
    expect(markup).toContain('按月');
  });

  it('renders Pagination correctly', () => {
    const markup = renderToStaticMarkup(
      <Pagination>
        <PaginationContent>
          <PaginationItem>
            <PaginationLink href="?page=1" isActive>1</PaginationLink>
          </PaginationItem>
        </PaginationContent>
      </Pagination>,
    );
    expect(markup).toContain('data-slot="pagination"');
    expect(markup).toContain('href="?page=1"');
  });

  it('renders OverflowTooltipText correctly', () => {
    const markup = renderToStaticMarkup(
      <OverflowTooltipText>超长门店门牌号与街道详细描述信息</OverflowTooltipText>,
    );
    expect(markup).toContain('超长门店门牌号与街道详细描述信息');
    expect(markup).toContain('truncate');
  });
});
