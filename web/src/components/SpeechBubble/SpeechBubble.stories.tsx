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
        text: "Hello, how is your day?",
        isOwn: false
    }
};

export const OwnMessage: Story = {
    args: {
        text: "Hello, how is your day?",
        isOwn: true
    }
};

export const MultipleMessages: Story = {
    args: {
        text: "Hello, how is your day?",
        isOwn: false
    },
    render: (args) => {
        return (
            <div className="flex flex-col gap-1">
                <SpeechBubble {...args} />
                <SpeechBubble {...args} text="My name is Jack" />
                <SpeechBubble
                    {...args}
                    text="Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec vel lacinia quam. Suspendisse interdum ultricies fringilla. Etiam condimentum magna quis risus vestibulum, sit amet eleifend odio tincidunt."
                />
            </div>
        );
    }
};
