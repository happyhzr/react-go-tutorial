import { Flex, Spinner, Stack, Text } from "@chakra-ui/react";
import TodoItem from "./TodoItem";
import { useQuery } from "@tanstack/react-query";
import { BASE_URL } from "../App";

export type Todo = {
    _id: number;
    body: string;
    completed: boolean;
}

const TodoList = () => {
    const query = useQuery<Todo[]>({
        queryKey: ['todos'],
        queryFn: async () => {
            try {
                const res = await fetch(`${BASE_URL}/todos`)
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
    })
    console.log(query)
    return (
        <>
            <Text fontSize={"4xl"} textTransform={"uppercase"} fontWeight={"bold"} textAlign={"center"} my={2} bgGradient='linear(to-l, #7928CA, #FF0080)' bgClip='text'>
                Today's Tasks
            </Text>
            {query.isLoading && (
                <Flex justifyContent={"center"} my={4}>
                    <Spinner size={"xl"} />
                </Flex>
            )}
            {!query.isLoading && query.data?.length === 0 && (
                <Stack alignItems={"center"} gap='3'>
                    <Text fontSize={"xl"} textAlign={"center"} color={"gray.500"}>
                        All tasks completed! 🤞
                    </Text>
                    <img src='/go.png' alt='Go logo' width={70} height={70} />
                </Stack>
            )}
            <Stack gap={3}>
                {query.data?.map((todo) => (
                    <TodoItem key={todo._id} todo={todo} />
                ))}
            </Stack>
        </>
    );
};
export default TodoList;