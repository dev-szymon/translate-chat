import type {Meta, StoryObj} from "@storybook/react";

import SpeechBubble from "./SpeechBubble";

const meta = {
    title: "atoms/SpeechBubble",
    component: SpeechBubble,
    parameters: {
        layout: "centered"
    }
} satisfies Meta<typeof SpeechBubble>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Default: Story = {
    args: {
        text: "test"
    }
};
