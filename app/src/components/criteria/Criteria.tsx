import { Box, Grid } from "@material-ui/core";
import React  from "react";
import Authors from "./Authors";
import Eras from "./Eras";
import Genres from "./Genres";
import * as models from "./models";

export default function Criteria (props: { onChange: (criteria: models.Criteria) => void }) {
  const [authors, setAuthors] = React.useState<models.Author[]>([]);
  const [genres, setGenres] = React.useState<models.Genre[]>([]);
  const [size, setSize] = React.useState<models.Size>(models.DefaultSize);
  const [era, setEra] = React.useState<models.Era>(models.DefaultEra);
  const [limit, setLimit] = React.useState<number>(10);

  function onChange() {
    props.onChange({
      authors: authors.map(x => x.id),
      genres: genres.map(x => x.id),
      minYear: era.minYear,
      maxYear: era.maxYear,
      minPages: size.minPages, 
      maxPages: size.maxPages, 
      limit: limit,
    });
  }

  return <Box>
    <h1>Criteria</h1>
    <Authors onChange={(x) => { setAuthors(x); onChange(); }}/>
    <Genres onChange={(x) => { setGenres(x); onChange(); }}/>
    <Grid container spacing={2}>
      <Grid item xs={4}>
        <Eras onChange={(x) => { setEra(x); onChange(); }}/>
      </Grid>
      <Grid item xs={4}>
      </Grid>
      <Grid item xs={4}>
      </Grid>
    </Grid>
  </Box>;
};
