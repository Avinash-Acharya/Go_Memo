import axios from "axios";
import Cards from "./_components/cards";
import CreateButton from "./_components/createbutton";

const URL =
  process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8000/api/go/users";

type dataType = {
  id: number;
  head: string;
  body: string;
  createdAt: string;
}[];

export const revalidate = 0;

async function getAllData() {
  try {
    // const data = await fetch(URL, { cache: "no-store" });
    // const allData: dataType = (await data.json()) as dataType;
    const response = await axios.get(URL);
    const allData: dataType = response.data as dataType;
    return allData.reverse();
  } catch (err) {
    console.log("Error while fetching all the notes : ", err);
  }
}

export default async function HomePage() {
  const allData: dataType | undefined = await getAllData();

  return (
    <div className="min-h-screen bg-gradient-to-b from-[#1a0b2e] to-[#000000] text-white">
      <header className="relative py-6 text-center">
        <h1 className="text-4xl">
          🧾
          <span>M</span>e<span>M</span>o
        </h1>
      </header>
      <div>
        <CreateButton />
      </div>
      <main>
        {allData ? (
          allData.map((data) => (
            <Cards
              key={data.id}
              id={data.id}
              head={data.head}
              body={data.body}
              createdAt={data.createdAt}
            />
          ))
        ) : (
          <h1 className="mx-auto mt-24 max-w-fit">Empty Repository</h1>
        )}
      </main>
    </div>
  );
}
