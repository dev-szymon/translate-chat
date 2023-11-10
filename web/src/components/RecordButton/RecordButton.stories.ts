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

export const Primary: Story = {
    args: {
        variant: "idle"
    }
};
