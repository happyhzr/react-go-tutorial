import { Button, Flex, Input, Spinner } from "@chakra-ui/react";
import { useMutation } from "@tanstack/react-query";
import React, { useState } from "react";
import { IoMdAdd } from "react-icons/io";
import { BASE_URL } from "../App";
import { useQueryClient } from "@tanstack/react-query";

const TodoForm = () => {
    const [newTodo, setNewTodo] = useState("");
    const queryClient = useQueryClient();
    const mutation = useMutation({
        mutationKey: ['createTodo'],
        mutationFn: async (e: React.FormEvent) => {
            e.preventDefault();
            try {
                const res = await fetch(`${BASE_URL}/todos`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ body: newTodo }),
                });
                const data = await res.json();
                if (!res.ok) {
                    throw new Error(data.message);
                }
                if (!data.success) {
                    throw new Error(data.message);
                }
                setNewTodo("");
                return data.data
            } catch (err) {
                console.log(err)
            }
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['todos'] });
        },
        onError: (err) => {
            alert(err.message)
        }
    })
    return (
        <form onSubmit={mutation.mutate}>
            <Flex gap={2}>
                <Input
                    type='text'
                    value={newTodo}
                    onChange={(e) => setNewTodo(e.target.value)}
                    ref={(input) => input && input.focus()}
                />
                <Button
                    mx={2}
                    type='submit'
                    _active={{
                        transform: "scale(.97)",
                    }}
                >
                    {mutation.isPending ? <Spinner size={"xs"} /> : <IoMdAdd size={30} />}
                </Button>
            </Flex>
        </form>
    );
};
export default TodoForm;