import React from "react";
import APICard from "./APICards";
import APIData from "../Dummy-Json/API-data";
import { Heading, Box } from "@chakra-ui/react";

const APIs = ({ query, category }) => {
  console.log(category);
  let apiCat = APIData.filter((apiData) => {
    if (category === null) {
      return apiData;
    } else if (apiData.categories.includes(category._id)) {
      return apiData;
    } else {
      return null;
    }
  });

  let apiCards = apiCat
    .filter((card) => {
      if (query === "") {
        return card;
      } else if (card.name.toLowerCase().includes(query.toLowerCase())) {
        return card;
      } else {
        return null;
      }
    })
    .map((item) => <APICard key={item.id} {...item} />);

  return (
    <Box>
      <Box id={"APIs"} h={1}></Box> {/*Remove for scooter*/}
      <Box my={50} mx={[50, 100, 150]}>
        <Heading as={"h3"} size={"lg"} mb={"50px"}>
          API's
        </Heading>
        {apiCards}
      </Box>
    </Box>
  );
};

export default APIs;