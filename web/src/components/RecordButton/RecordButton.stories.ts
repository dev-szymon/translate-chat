import type {Meta, StoryObj} from "@storybook/react";

import RecordButton from "./RecordButton";

const meta = {
    title: "atoms/RecordButton",
    component: RecordButton,
    parameters: {
        layout: "centered"
    },
    argTypes: {
        variant: {control: "radio", options: ["idle", "recording", "pending"]}
    }
} satisfies Meta<typeof RecordButton>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Idle: Story = {
    args: {
        variant: "idle"
    }
};

export const Recording: Story = {
    args: {
        variant: "recording"
    }
};

export const Pending: Story = {
    args: {
        variant: "pending"
    }
};
