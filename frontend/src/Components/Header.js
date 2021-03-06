import {
  Box,
  Flex,
  Text,
  Heading,
  VStack,
  Input,
  InputGroup,
  InputRightElement,
  Center,
} from "@chakra-ui/react";
import { Link } from "react-scroll";
import { React } from "react";
import "../theme/styles.css";

const Header = ({ query, setQuery }) => {
  return (
    <>
      <Box h={"90vh"} mx={50}>
        <Flex h={"80%"} alignItems={"center"} justifyContent={"center"}>
          <VStack spacing={10}>
            <Heading fontSize={["4rem", "4.5rem", "5.25rem", "8rem"]}>
              <span id="the">The</span>
              <span id="modern">Modern</span>
              <br />
              <span id="api">API</span>
              <span id="platform">Platform</span>
            </Heading>
            <Text
              textAlign={"center"}
              fontSize={["1rem", null, "1.25rem", "1.5rem"]}
              color={"#222222"}
            >
              Find, connect to, &amp; manage thousands of API's
            </Text>
          </VStack>
        </Flex>
        <Center>
          <InputGroup w={["80vw", "70vw", "60vw", "50vw", "40vw"]}>
            <Input
              fontSize={["1rem", null, "1.25rem", "1.5rem"]}
              onChange={(event) => setQuery(event.target.value)}
              isInvalid
              errorBorderColor="#222222"
              noOfLines={1}
              placeholder="Search for an API here"
              value={query}
            ></Input>

            <InputRightElement mr={"1px"} w={"auto"} my={"auto"}>
              <Link
                to="APIs"
                spy={true}
                smooth={true}
                offset={-50}
                duration={500}
              >
                <Box
                  px={[1, 2, 3, 4, 5]}
                  py={["7px", null, "4px", "1px"]}
                  borderEndRadius={"5px"}
                  as="button"
                  fontSize={["1rem", null, "1.25rem", "1.5rem"]}
                  bg="#ffffff"
                  color={"#222222"}
                  _hover={{
                    bg: "#222222",
                    color: "#ffffff",
                  }}
                >
                  Search
                </Box>
              </Link>
            </InputRightElement>
          </InputGroup>
        </Center>
      </Box>
    </>
  );
};

export default Header;
