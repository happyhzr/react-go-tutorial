import { Badge, Box, Flex, Text, Spinner } from "@chakra-ui/react";
import { FaCheckCircle } from "react-icons/fa";
import { MdDelete } from "react-icons/md";
import type { Todo } from "./TodoList";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { BASE_URL } from "../App";

const TodoItem = ({ todo }: { todo: Todo }) => {
    const queryClient = useQueryClient();
    const mutation = useMutation({
        mutationKey: ['updateTodo'],
        mutationFn: async () => {
            if (todo.completed) {
                return alert("Task already completed");
            }
            try {
                const res = await fetch(`${BASE_URL}/todos/${todo._id}`, { method: 'PATCH' })
                const data = await res.json()
                if (!res.ok) {
                    throw new Error(data.message)
                }
                if (!data.success) {
                    throw new Error(data.message)
                }
                return data.data
            } catch (err) {
                console.log(err)
            }
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['todos'] });
        },
    })
    const deleteMutation = useMutation({
        mutationKey: ['deleteTodo'],
        mutationFn: async () => {
            try {
                const res = await fetch(`${BASE_URL}/todos/${todo._id}`, { method: 'DELETE' })
                const data = await res.json()
                if (!res.ok) {
                    throw new Error(data.message)
                }
                if (!data.success) {
                    throw new Error(data.message)
                }
                return data.data
            } catch (err) {
                console.log(err)
            }
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['todos'] });
        },
    })
    return (
        <Flex gap={2} alignItems={"center"} >
            <Flex
                flex={1}
                alignItems={"center"}
                border={"1px"}
                borderColor={"gray.600"}
                p={2}
                borderRadius={"lg"}
                justifyContent={"space-between"}
            >
                <Text
                    color={todo.completed ? "green.200" : "yellow.100"}
                    textDecoration={todo.completed ? "line-through" : "none"}
                >
                    {todo.body}
                </Text>
                {todo.completed && (
                    <Badge ml='1' colorScheme='green'>
                        Done
                    </Badge>
                )}
                {!todo.completed && (
                    <Badge ml='1' colorScheme='yellow'>
                        In Progress
                    </Badge>
                )}
            </Flex>
            <Flex gap={2} alignItems={"center"}>
                <Box color={"green.500"} cursor={"pointer"} onClick={() => mutation.mutate()}>
                    {!mutation.isPending && <FaCheckCircle size={20} />}
                    {mutation.isPending && <Spinner size={"sm"} />}
                </Box>
                <Box color={"red.500"} cursor={"pointer"} onClick={() => deleteMutation.mutate()}>
                    {!deleteMutation.isPending && <MdDelete size={25} />}
                    {deleteMutation.isPending && <Spinner size={"sm"} />}
                </Box>
            </Flex>
        </Flex >
    );
};
export default TodoItem;