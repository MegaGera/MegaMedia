'use client';

import Link from 'next/link';

import {Card, CardBody, CardFooter, Image} from "@nextui-org/react";

export default function App() {
  const list = [
    {
      title: "MegaGera",
      img: "/images/sarita.jpg",
      href: "/dashboard/megagera"
    },
    {
      title: "MegaGoal",
      img: "/images/fruit-1.jpg",
      href: "/dashboard/megagoal"
    }
  ];

  return (
    <div className="gap-2 grid grid-cols-2 sm:grid-cols-4">
      {list.map((item, index) => (
        <Link
            key={item.title}
            href={item.href}
            className='w-full'
          >
          <Card className='w-[90%] m-auto' key={index} isPressable shadow="sm" onPress={() => console.log("item pressed")}>
            <CardBody className="overflow-visible p-0">
              <Image
                alt={item.title}
                className="w-full object-cover h-[200px]"
                radius="lg"
                shadow="sm"
                src={item.img}
                width="100%"
              />
            </CardBody>
            <CardFooter className="text-small justify-between">
              <b>{item.title}</b>
              <p className="text-default-500">Access</p>
            </CardFooter>
          </Card>
        </Link>
      ))}
    </div>
  );
}
