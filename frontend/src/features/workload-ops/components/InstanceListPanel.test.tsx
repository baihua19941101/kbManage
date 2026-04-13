import '@testing-library/jest-dom/vitest';
import { fireEvent, render, screen } from '@testing-library/react';
import { InstanceListPanel } from '@/features/workload-ops/components/InstanceListPanel';
import { installAntdDomShims } from '@/test/installAntdDomShims';

describe('InstanceListPanel', () => {
  beforeAll(() => {
    installAntdDomShims();
  });

  it('renders instances and triggers terminal action', async () => {
    const onOpenTerminal = vi.fn();
    render(
      <InstanceListPanel
        items={[
          {
            podName: 'demo-api-pod-0',
            containerName: 'app',
            phase: 'Running',
            ready: true,
            terminalAvailable: true
          }
        ]}
        onOpenTerminal={onOpenTerminal}
      />
    );

    expect(await screen.findByText('demo-api-pod-0')).toBeInTheDocument();
    fireEvent.click(screen.getByRole('button', { name: /终\s*端/ }));
    expect(onOpenTerminal).toHaveBeenCalledTimes(1);
  });
});
