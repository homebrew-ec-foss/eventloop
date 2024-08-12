import "./globals.css";
import { isPrime } from "mathjs";
import React, { useState, useEffect } from "react";
import { NavigationMenu, NavigationMenuItem, NavigationMenuList, NavigationMenuLink } from "@radix-ui/react-navigation-menu"
import Link from "next/link"
import { navigationMenuTriggerStyle } from "@/components/ui/navigation-menu"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { Badge } from "@/components/ui/badge";

export default function Ping() {
    const [response, setResponse] = useState(null);

    useEffect(() => {
        async function fetchData() {
            try {
                const res = await fetch('http://localhost:8080/ping', {
                    method: 'GET',
                });
                if (res.ok) {
                    const text = await res.text();
                    setResponse(text);
                } else {
                    throw new Error('Failed to fetch');
                }
            } catch (error) {
                // console.log("Error while fetching ping response: ", error);
                setResponse(null, error);
            }
        }
        // setInterval(fetchData, 10000);
        fetchData();
    }, []);

    return (
        <main className="flex min-h-screen flex-col p-5 md:p-24 gap-4">
            <NavigationMenu>
                <NavigationMenuList>
                    <NavigationMenuItem>
                        <Link href="/" legacyBehavior passHref>
                            <NavigationMenuLink className={navigationMenuTriggerStyle()}>
                                Home
                            </NavigationMenuLink>
                        </Link>
                        <Link href="/create" legacyBehavior passHref>
                            <NavigationMenuLink className={navigationMenuTriggerStyle()}>
                                Create
                            </NavigationMenuLink>
                        </Link>
                        <Link href="/ping" legacyBehavior passHref>
                            <NavigationMenuLink className={navigationMenuTriggerStyle()}>
                                Ping
                            </NavigationMenuLink>
                        </Link>
                    </NavigationMenuItem>
                </NavigationMenuList>
            </NavigationMenu>

            <Card>
                <CardHeader>
                    <CardTitle>Server Status</CardTitle>
                    <CardDescription>Eventloop GO backend status</CardDescription>
                </CardHeader>
                <CardContent>
                    {
                        response ? (
                            // <p className="text-2xl font-mono">{response}</p>
                            <Badge>Alive</Badge>
                        ): (
                            <Badge variant="destructive">Not responsive</Badge>
                        )
                    }
                </CardContent>
            </Card>

            {/* <p className={`text-3xl font-regular underline`}>{response}</p> */}
        </main>
    );
}
