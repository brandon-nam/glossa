import {
  Avatar,
  AvatarFallback,
  AvatarImage,
} from '@/components/ui/avatar';
import type { ColumnDef } from '@/components/ui/shadcn-io/table';
import {
  TableBody,
  TableCell,
  TableColumnHeader,
  TableHead,
  TableHeader,
  TableHeaderGroup,
  TableProvider,
  TableRow,
} from '@/components/ui/shadcn-io/table';
import { ChevronRightIcon } from 'lucide-react';
// Fixed data to prevent hydration mismatches
const statuses = [
  { id: 'status-1', name: 'Planned', color: '#6B7280' },
  { id: 'status-2', name: 'In Progress', color: '#F59E0B' },
  { id: 'status-3', name: 'Done', color: '#10B981' },
];
const users = [
  { id: 'user-1', name: 'John Doe', image: 'https://github.com/shadcn.png' },
  { id: 'user-2', name: 'Jane Smith', image: 'https://github.com/vercel.png' },
  { id: 'user-3', name: 'Bob Johnson', image: 'https://github.com/nextjs.png' },
  { id: 'user-4', name: 'Alice Brown', image: 'https://github.com/tailwindcss.png' },
];
const exampleGroups = [
  { id: 'group-1', name: 'Frontend Team' },
  { id: 'group-2', name: 'Backend Team' },
  { id: 'group-3', name: 'Design Team' },
  { id: 'group-4', name: 'Product Team' },
  { id: 'group-5', name: 'Marketing Team' },
  { id: 'group-6', name: 'Sales Team' },
];
const exampleProducts = [
  { id: 'product-1', name: 'Web Platform' },
  { id: 'product-2', name: 'Mobile App' },
  { id: 'product-3', name: 'API Gateway' },
  { id: 'product-4', name: 'Analytics Dashboard' },
];
const exampleInitiatives = [
  { id: 'initiative-1', name: 'Q1 Product Launch' },
  { id: 'initiative-2', name: 'Performance Optimization' },
];
const exampleReleases = [
  { id: 'release-1', name: 'Version 2.0' },
  { id: 'release-2', name: 'Version 2.1' },
  { id: 'release-3', name: 'Version 3.0' },
];
const exampleFeatures = [
  {
    id: 'feature-1',
    name: 'User Authentication System',
    startAt: new Date('2024-01-15'),
    endAt: new Date('2024-02-15'),
    status: statuses[1],
    owner: users[0],
    group: exampleGroups[0],
    product: exampleProducts[0],
    initiative: exampleInitiatives[0],
    release: exampleReleases[0],
  },
  {
    id: 'feature-2',
    name: 'Dashboard Analytics',
    startAt: new Date('2024-02-01'),
    endAt: new Date('2024-03-01'),
    status: statuses[2],
    owner: users[1],
    group: exampleGroups[1],
    product: exampleProducts[1],
    initiative: exampleInitiatives[1],
    release: exampleReleases[1],
  },
  {
    id: 'feature-3',
    name: 'API Rate Limiting',
    startAt: new Date('2024-01-20'),
    endAt: new Date('2024-02-20'),
    status: statuses[0],
    owner: users[2],
    group: exampleGroups[2],
    product: exampleProducts[2],
    initiative: exampleInitiatives[0],
    release: exampleReleases[2],
  },
  {
    id: 'feature-4',
    name: 'Mobile Push Notifications',
    startAt: new Date('2024-03-01'),
    endAt: new Date('2024-04-01'),
    status: statuses[1],
    owner: users[3],
    group: exampleGroups[3],
    product: exampleProducts[1],
    initiative: exampleInitiatives[1],
    release: exampleReleases[0],
  },
  {
    id: 'feature-5',
    name: 'Real-time Chat System',
    startAt: new Date('2024-02-15'),
    endAt: new Date('2024-03-15'),
    status: statuses[2],
    owner: users[0],
    group: exampleGroups[0],
    product: exampleProducts[0],
    initiative: exampleInitiatives[0],
    release: exampleReleases[1],
  },
];
const Example = () => {
  const columns: ColumnDef<(typeof exampleFeatures)[number]>[] = [
    {
      accessorKey: 'name',
      header: ({ column }) => (
        <TableColumnHeader column={column} title="Name" />
      ),
      cell: ({ row }) => (
        <div className="flex items-center gap-2">
          <div className="relative">
            <Avatar className="size-6">
              <AvatarImage src={row.original.owner.image} />
              <AvatarFallback>
                {row.original.owner.name.split(' ').map(n => n[0]).join('').slice(0, 2)}
              </AvatarFallback>
            </Avatar>
            <div
              className="absolute right-0 bottom-0 h-2 w-2 rounded-full ring-2 ring-background"
              style={{
                backgroundColor: row.original.status.color,
              }}
            />
          </div>
          <div>
            <span className="font-medium">{row.original.name}</span>
            <div className="flex items-center gap-1 text-muted-foreground text-xs">
              <span>{row.original.product.name}</span>
              <ChevronRightIcon size={12} />
              <span>{row.original.group.name}</span>
            </div>
          </div>
        </div>
      ),
    },
    {
      accessorKey: 'startAt',
      header: ({ column }) => (
        <TableColumnHeader column={column} title="Start At" />
      ),
      cell: ({ row }) => {
        const date = new Date(row.original.startAt);
        return date.toLocaleDateString('en-US', {
          year: 'numeric',
          month: 'short',
          day: 'numeric',
        });
      },
    },
    {
      accessorKey: 'endAt',
      header: ({ column }) => (
        <TableColumnHeader column={column} title="End At" />
      ),
      cell: ({ row }) => {
        const date = new Date(row.original.endAt);
        return date.toLocaleDateString('en-US', {
          year: 'numeric',
          month: 'short',
          day: 'numeric',
        });
      },
    },
    {
      id: 'release',
      accessorFn: (row) => row.release.id,
      header: ({ column }) => (
        <TableColumnHeader column={column} title="Release" />
      ),
      cell: ({ row }) => row.original.release.name,
    },
  ];
  return (
    <TableProvider columns={columns} data={exampleFeatures}>
      <TableHeader>
        {({ headerGroup }) => (
          <TableHeaderGroup headerGroup={headerGroup} key={headerGroup.id}>
            {({ header }) => <TableHead header={header} key={header.id} />}
          </TableHeaderGroup>
        )}
      </TableHeader>
      <TableBody>
        {({ row }) => (
          <TableRow key={row.id} row={row}>
            {({ cell }) => <TableCell cell={cell} key={cell.id} />}
          </TableRow>
        )}
      </TableBody>
    </TableProvider>
  );
};
export default Example;
